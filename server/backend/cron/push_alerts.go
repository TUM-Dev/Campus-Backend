package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/TUM-Dev/Campus-Backend/server/internal/fcm"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type trackedServiceDef struct {
	key     string
	labelDE string
	labelEN string
	pattern *regexp.Regexp
}

var trackedITServices = []trackedServiceDef{
	{key: "tumonline", labelDE: "TUMonline", labelEN: "TUMonline", pattern: tumonlineService},
	{key: "moodle", labelDE: "Moodle", labelEN: "Moodle", pattern: regexp.MustCompile(`(?i)\bmoodle\b`)},
	{key: "eduroam", labelDE: "eduroam / WLAN", labelEN: "eduroam / Wi-Fi", pattern: regexp.MustCompile(`(?i)\b(eduroam|wlan|wi-?fi)\b`)},
}

type parsedServiceStatus struct {
	Key       string
	Status    string
	MessageDE string
	MessageEN string
}

// parseTrackedServiceStatuses extracts active outage/maintenance states for TUMonline, Moodle, and eduroam
// from the TUM IT status page HTML.
func parseTrackedServiceStatuses(r io.Reader) ([]parsedServiceStatus, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	var current *goquery.Selection
	var future, past bool
	doc.Find("h2").Each(func(_ int, heading *goquery.Selection) {
		switch strings.Join(strings.Fields(heading.Text()), " ") {
		case "Aktuelle Meldungen", "Current notifications":
			current = heading.NextUntil("h2")
		case "Angekündigte Meldungen", "Future notifications":
			future = true
		case "Beendete Meldungen ↑", "Past notifications ↑":
			past = true
		}
	})
	if current == nil || !future || !past || current.Find(".table-headers").Length() != 1 {
		return nil, fmt.Errorf("unrecognized TUM status page structure")
	}

	byKey := make(map[string]*parsedServiceStatus, len(trackedITServices))
	for _, svc := range trackedITServices {
		byKey[svc.key] = &parsedServiceStatus{Key: svc.key, Status: model.ServiceStatusOK}
	}

	current.Find("a.table-entry").Each(func(_ int, entry *goquery.Selection) {
		serviceText := strings.TrimSpace(entry.Find(".entry-service .mobile-bold").Text())
		title := strings.TrimSpace(entry.Find(".entry-problemtype").Text())
		if serviceText == "" || title == "" {
			return
		}
		isMaint := maintenanceNotice.MatchString(title)
		for _, def := range trackedITServices {
			if !def.pattern.MatchString(serviceText) {
				continue
			}
			st := byKey[def.key]
			if isMaint {
				st.Status = model.ServiceStatusMaintenance
				st.MessageDE = fmt.Sprintf("%s wird derzeit gewartet (%s).", def.labelDE, title)
				st.MessageEN = fmt.Sprintf("%s is currently under maintenance (%s).", def.labelEN, title)
			} else if st.Status != model.ServiceStatusMaintenance {
				st.Status = model.ServiceStatusDegraded
				st.MessageDE = fmt.Sprintf("%s ist derzeit beeinträchtigt (%s).", def.labelDE, title)
				st.MessageEN = fmt.Sprintf("%s is currently degraded (%s).", def.labelEN, title)
			}
		}
	})

	out := make([]parsedServiceStatus, 0, len(trackedITServices))
	for _, def := range trackedITServices {
		out = append(out, *byKey[def.key])
	}
	return out, nil
}

type pushRecipient struct {
	DeviceID             string `gorm:"column:device_id"`
	FCMToken             string `gorm:"column:fcm_token"`
	CanteenIDs           string `gorm:"column:canteen_ids"`
	FavoriteDishes       string `gorm:"column:favorite_dishes"`
	DietaryLabels        string `gorm:"column:dietary_labels"`
	NotifyFavoriteDishes bool   `gorm:"column:notify_favorite_dishes"`
	NotifyITStatus       bool   `gorm:"column:notify_it_status"`
	Language             string `gorm:"column:language"`
}

func decodeJSONList(raw string) []string {
	if raw == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

// syncAndNotifyServiceStatuses updates service_status_states and sends FCM notifications
// whenever a tracked IT service transitions between ok and degraded/maintenance.
func syncAndNotifyServiceStatuses(ctx context.Context, db *gorm.DB, sender fcm.Sender, statuses []parsedServiceStatus) error {
	var transitions []struct {
		prev   string
		latest parsedServiceStatus
	}
	for _, latest := range statuses {
		var prev model.ServiceStatusState
		err := db.WithContext(ctx).First(&prev, "service_key = ?", latest.Key).Error
		prevStatus := model.ServiceStatusOK
		if err == nil {
			prevStatus = prev.Status
		}
		state := model.ServiceStatusState{
			ServiceKey: latest.Key,
			Status:     latest.Status,
			MessageDE:  latest.MessageDE,
			MessageEN:  latest.MessageEN,
			UpdatedAt:  time.Now(),
		}
		if err := db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "service_key"}},
			DoUpdates: clause.AssignmentColumns([]string{"status", "message_de", "message_en", "updated_at"}),
		}).Create(&state).Error; err != nil {
			return err
		}
		if prevStatus != latest.Status {
			transitions = append(transitions, struct {
				prev   string
				latest parsedServiceStatus
			}{prev: prevStatus, latest: latest})
		}
	}
	if len(transitions) == 0 || sender == nil {
		return nil
	}

	var recipients []pushRecipient
	if err := db.WithContext(ctx).
		Table("push_preferences p").
		Joins("JOIN push_devices d ON d.device_id = p.device_id").
		Select("p.device_id, d.fcm_token, p.language").
		Where("p.notify_it_status = 1 AND d.fcm_token <> ''").
		Scan(&recipients).Error; err != nil {
		return err
	}

	for _, tr := range transitions {
		for _, r := range recipients {
			title, body := buildServiceStatusNotification(tr.latest, tr.prev, r.Language)
			if err := sender.Send(ctx, fcm.Message{
				Token: r.FCMToken,
				Title: title,
				Body:  body,
				Data: map[string]string{
					"type":    "it_status",
					"service": tr.latest.Key,
					"status":  tr.latest.Status,
				},
			}); err != nil {
				log.WithError(err).WithFields(log.Fields{"device_id": r.DeviceID, "service": tr.latest.Key}).Warn("failed to send IT status push")
			}
		}
	}
	return nil
}

func buildServiceStatusNotification(latest parsedServiceStatus, prevStatus, lang string) (string, string) {
	label := latest.Key
	for _, def := range trackedITServices {
		if def.key == latest.Key {
			if lang == "de" {
				label = def.labelDE
			} else {
				label = def.labelEN
			}
			break
		}
	}
	if latest.Status == model.ServiceStatusOK && prevStatus != model.ServiceStatusOK {
		if lang == "de" {
			return label + " wieder verfügbar", label + " läuft wieder normal."
		}
		return label + " recovered", label + " is back to normal."
	}
	if lang == "de" {
		body := latest.MessageDE
		if body == "" {
			body = label + " meldet eine Störung."
		}
		return label + " Status", body
	}
	body := latest.MessageEN
	if body == "" {
		body = label + " is experiencing an issue."
	}
	return label + " Status", body
}

type scheduledMenuDish struct {
	DishID      int64  `gorm:"column:dishID"`
	DishName    string `gorm:"column:dishName"`
	DishLabels  string `gorm:"column:labels"`
	CanteenName string `gorm:"column:canteenName"`
	Day         int32  `gorm:"column:day"`
}

// notifyFavoriteDishAlerts sends deduplicated FCM notifications for upcoming menu dishes
// matching a device's favorite_dishes keywords, dietary_labels, and monitored canteens.
func notifyFavoriteDishAlerts(ctx context.Context, db *gorm.DB, sender fcm.Sender, year, week int) error {
	if sender == nil {
		return nil
	}
	var recipients []pushRecipient
	if err := db.WithContext(ctx).
		Table("push_preferences p").
		Joins("JOIN push_devices d ON d.device_id = p.device_id").
		Select("p.device_id, d.fcm_token, p.canteen_ids, p.favorite_dishes, p.dietary_labels, p.language").
		Where("p.notify_favorite_dishes = 1 AND d.fcm_token <> ''").
		Scan(&recipients).Error; err != nil {
		return err
	}
	if len(recipients) == 0 {
		return nil
	}

	var menu []scheduledMenuDish
	if err := db.WithContext(ctx).
		Table("dishes_of_the_weeks w").
		Joins("JOIN dishes d ON d.dish = w.dishID").
		Joins("JOIN canteens c ON c.cafeteria = d.cafeteriaID").
		Select("DISTINCT d.dish AS dishID, d.name AS dishName, COALESCE(d.labels, '') AS labels, c.name AS canteenName, w.day AS day").
		Where("w.year = ? AND w.week = ?", year, week).
		Scan(&menu).Error; err != nil {
		return err
	}
	if len(menu) == 0 {
		return nil
	}

	for _, r := range recipients {
		favorites := decodeJSONList(r.FavoriteDishes)
		if len(favorites) == 0 {
			continue
		}
		allowedCanteens := make(map[string]bool)
		for _, c := range decodeJSONList(r.CanteenIDs) {
			allowedCanteens[model.CanteenEatAPIID(c)] = true
		}
		requiredLabels := decodeJSONList(r.DietaryLabels)

		for _, item := range menu {
			canteenEatID := model.CanteenEatAPIID(item.CanteenName)
			if len(allowedCanteens) > 0 && !allowedCanteens[canteenEatID] {
				continue
			}
			if !matchesFavoriteDish(item.DishName, favorites) {
				continue
			}
			if !matchesAllDietaryLabels(item.DishLabels, requiredLabels) {
				continue
			}

			dedupKey := fmt.Sprintf("fav-dish:%d:%d:%d:%d", item.DishID, year, week, item.Day)
			logEntry := model.PushDeliveryLog{
				DeviceID: r.DeviceID,
				DedupKey: dedupKey,
				SentAt:   time.Now(),
			}
			res := db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&logEntry)
			if res.Error != nil || res.RowsAffected == 0 {
				continue
			}

			title, body := buildFavoriteDishNotification(item.DishName, canteenEatID, r.Language)
			if err := sender.Send(ctx, fcm.Message{
				Token: r.FCMToken,
				Title: title,
				Body:  body,
				Data: map[string]string{
					"type":       "favorite_dish",
					"canteen_id": canteenEatID,
					"dish":       item.DishName,
				},
			}); err != nil {
				log.WithError(err).WithFields(log.Fields{"device_id": r.DeviceID, "dish": item.DishName}).Warn("failed to send favorite dish push")
			}
		}
	}
	return nil
}

func matchesFavoriteDish(dishName string, favorites []string) bool {
	lower := strings.ToLower(dishName)
	for _, fav := range favorites {
		f := strings.ToLower(strings.TrimSpace(fav))
		if f != "" && strings.Contains(lower, f) {
			return true
		}
	}
	return false
}

func matchesAllDietaryLabels(rawLabels string, required []string) bool {
	if len(required) == 0 {
		return true
	}
	have := make(map[string]bool)
	for _, part := range strings.Split(rawLabels, ",") {
		if t := strings.ToUpper(strings.TrimSpace(part)); t != "" {
			have[t] = true
		}
	}
	for _, req := range required {
		r := strings.ToUpper(strings.TrimSpace(req))
		if r != "" && !have[r] {
			return false
		}
	}
	return true
}

func buildFavoriteDishNotification(dishName, canteenID, lang string) (string, string) {
	if lang == "de" {
		return "Lieblingsgericht auf dem Speiseplan!", fmt.Sprintf("%s gibt es bei %s.", dishName, canteenID)
	}
	return "Favorite dish on the menu!", fmt.Sprintf("%s is served at %s.", dishName, canteenID)
}
