package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/guregu/null"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	ratingImagesEnabledEnv      = "RATING_IMAGES_ENABLED"
	discordModerationWebhookEnv = "DISCORD_MODERATION_WEBHOOK_URL"
	maxImageDecodeDimension     = 8000
	discordWebhookTimeout       = 5 * time.Second
)

// ratingImagesEnabled returns true only when RATING_IMAGES_ENABLED=true is explicitly set.
// Rating image uploads are disabled by default to prevent unscreened user uploads.
func ratingImagesEnabled() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv(ratingImagesEnabledEnv)), "true")
}

func checkDeviceNotBanned(tx *gorm.DB, deviceHash null.String) error {
	if !deviceHash.Valid {
		return nil
	}
	var count int64
	if err := tx.Model(&model.BannedRatingDevice{}).
		Where("device_hash = ?", deviceHash.String).
		Count(&count).Error; err != nil {
		log.WithError(err).Error("while checking banned rating devices")
		return status.Error(codes.Internal, "Could not verify device status.")
	}
	if count > 0 {
		return status.Error(codes.PermissionDenied, "This device is restricted from submitting ratings or reports.")
	}
	return nil
}

// ReportedDishRating is a summary view of a reported or hidden dish rating, designed to be
// rendered directly by the opt-in /admin/ moderation dashboard.
type ReportedDishRating struct {
	RatingID    int64     `gorm:"column:rating_id"`
	DishName    string    `gorm:"column:dish_name"`
	CanteenName string    `gorm:"column:canteen_name"`
	Points      int32     `gorm:"column:points"`
	Comment     string    `gorm:"column:comment"`
	Status      string    `gorm:"column:status"`
	DeviceHash  string    `gorm:"column:device_hash"`
	Timestamp   time.Time `gorm:"column:timestamp"`
	ReportCount int64     `gorm:"column:report_count"`
	Reasons     string    `gorm:"column:reasons"`
}

// SetDishRatingStatus hides or restores a dish rating without altering its original visit timestamp.
func (s *CampusServer) SetDishRatingStatus(ctx context.Context, ratingID int64, newStatus string) error {
	if ratingID <= 0 {
		return fmt.Errorf("invalid rating_id")
	}
	if newStatus != model.RatingStatusVisible && newStatus != model.RatingStatusHidden {
		return fmt.Errorf("invalid status %q", newStatus)
	}
	res := s.db.WithContext(ctx).Model(&model.DishRating{}).
		Where("dishRating = ?", ratingID).
		Updates(map[string]any{
			"status":    newStatus,
			"timestamp": gorm.Expr("timestamp"),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// BanRatingDevice adds a device_hash to banned_rating_devices and optionally hides all
// existing dish and canteen ratings submitted by that device.
func (s *CampusServer) BanRatingDevice(ctx context.Context, deviceHash, reason string, hideExisting bool) error {
	deviceHash = strings.TrimSpace(deviceHash)
	if deviceHash == "" {
		return fmt.Errorf("device_hash is required")
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		entry := model.BannedRatingDevice{
			DeviceHash: deviceHash,
			Reason:     strings.TrimSpace(reason),
			BannedAt:   time.Now(),
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "device_hash"}},
			DoUpdates: clause.AssignmentColumns([]string{"reason", "banned_at"}),
		}).Create(&entry).Error; err != nil {
			return err
		}
		if !hideExisting {
			return nil
		}
		if err := tx.Model(&model.DishRating{}).
			Where("device_hash = ?", deviceHash).
			Updates(map[string]any{
				"status":    model.RatingStatusHidden,
				"timestamp": gorm.Expr("timestamp"),
			}).Error; err != nil {
			return err
		}
		return tx.Model(&model.CanteenRating{}).
			Where("device_hash = ?", deviceHash).
			Updates(map[string]any{
				"status":    model.RatingStatusHidden,
				"timestamp": gorm.Expr("timestamp"),
			}).Error
	})
}

// UnbanRatingDevice removes a device_hash from banned_rating_devices.
func (s *CampusServer) UnbanRatingDevice(ctx context.Context, deviceHash string) error {
	deviceHash = strings.TrimSpace(deviceHash)
	if deviceHash == "" {
		return fmt.Errorf("device_hash is required")
	}
	return s.db.WithContext(ctx).Where("device_hash = ?", deviceHash).Delete(&model.BannedRatingDevice{}).Error
}

// ListReportedDishRatings returns dish ratings that have at least one user report or are hidden,
// ordered by report count descending.
func (s *CampusServer) ListReportedDishRatings(ctx context.Context, limit, offset int) ([]ReportedDishRating, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	var rows []ReportedDishRating
	err := s.db.WithContext(ctx).
		Table("dish_ratings r").
		Joins("JOIN dishes d ON d.dish = r.dishID").
		Joins("JOIN canteens c ON c.cafeteria = d.cafeteriaID").
		Joins("LEFT JOIN dish_rating_reports rep ON rep.rating_id = r.dishRating").
		Select(`r.dishRating AS rating_id,
			d.name AS dish_name,
			c.name AS canteen_name,
			r.points AS points,
			COALESCE(r.comment, '') AS comment,
			r.status AS status,
			COALESCE(r.device_hash, '') AS device_hash,
			r.timestamp AS timestamp,
			COUNT(rep.id) AS report_count,
			COALESCE(GROUP_CONCAT(NULLIF(rep.reason, '') SEPARATOR '; '), '') AS reasons`).
		Where("r.status = ? OR rep.id IS NOT NULL", model.RatingStatusHidden).
		Group("r.dishRating, d.name, c.name, r.points, r.comment, r.status, r.device_hash, r.timestamp").
		Order("report_count DESC, r.timestamp DESC").
		Limit(limit).
		Offset(offset).
		Scan(&rows).Error
	return rows, err
}

type moderationReportAlert struct {
	RatingID     int64
	DishName     string
	CanteenID    string
	Points       int32
	Comment      string
	Reason       string
	TotalReports int64
	AutoHidden   bool
	DeviceHash   string
}

type discordWebhookPayload struct {
	Embeds []discordEmbed `json:"embeds"`
}

type discordEmbed struct {
	Title  string              `json:"title"`
	Color  int                 `json:"color"`
	Fields []discordEmbedField `json:"fields"`
}

type discordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

func (s *CampusServer) moderationWebhookTarget() string {
	if s.moderationWebhookURL != "" {
		return s.moderationWebhookURL
	}
	return strings.TrimSpace(os.Getenv(discordModerationWebhookEnv))
}

// notifyModerationDiscordAsync posts a text-only moderation embed to Discord without blocking the RPC.
// It never includes user-uploaded image URLs so reviewers are never exposed to unscreened images.
func (s *CampusServer) notifyModerationDiscordAsync(alert moderationReportAlert) {
	webhookURL := s.moderationWebhookTarget()
	if webhookURL == "" {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), discordWebhookTimeout)
		defer cancel()
		if err := s.sendDiscordModerationAlert(ctx, webhookURL, alert); err != nil {
			log.WithError(err).WithField("rating_id", alert.RatingID).Warn("failed to deliver Discord moderation alert")
		}
	}()
}

func (s *CampusServer) sendDiscordModerationAlert(ctx context.Context, webhookURL string, alert moderationReportAlert) error {
	title := fmt.Sprintf("Dish rating #%d reported (%d total)", alert.RatingID, alert.TotalReports)
	color := 0xF59E0B // amber
	if alert.AutoHidden {
		title = fmt.Sprintf("Dish rating #%d auto-hidden (%d reports)", alert.RatingID, alert.TotalReports)
		color = 0xEF4444 // red
	}
	comment := alert.Comment
	if comment == "" {
		comment = "(no comment)"
	}
	reason := alert.Reason
	if reason == "" {
		reason = "(none provided)"
	}
	devicePrefix := "(anonymous)"
	if len(alert.DeviceHash) >= 12 {
		devicePrefix = alert.DeviceHash[:12] + "…"
	} else if alert.DeviceHash != "" {
		devicePrefix = alert.DeviceHash
	}

	payload := discordWebhookPayload{
		Embeds: []discordEmbed{{
			Title: title,
			Color: color,
			Fields: []discordEmbedField{
				{Name: "Canteen", Value: alert.CanteenID, Inline: true},
				{Name: "Dish", Value: alert.DishName, Inline: true},
				{Name: "Score", Value: fmt.Sprintf("%d/5", alert.Points), Inline: true},
				{Name: "Comment", Value: comment},
				{Name: "Report Reason", Value: reason},
				{Name: "Author Device Hash", Value: "`" + devicePrefix + "`", Inline: true},
			},
		}},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := s.moderationHTTPClient
	if client == nil {
		client = &http.Client{Timeout: discordWebhookTimeout}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return errors.New("discord webhook returned " + resp.Status + ": " + string(respBody))
	}
	return nil
}
