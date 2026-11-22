package backend

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	maxFavoriteDishKeywords = 50
	maxDietaryLabels        = 20
	maxCanteenPrefs         = 30
)

func normalizeStringList(items []string, transform func(string) string) []string {
	seen := make(map[string]bool, len(items))
	out := make([]string, 0, len(items))
	for _, raw := range items {
		val := strings.TrimSpace(raw)
		if val == "" {
			continue
		}
		if transform != nil {
			val = transform(val)
		}
		key := strings.ToLower(val)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, val)
	}
	return out
}

func encodeStringSliceJSON(items []string) string {
	if len(items) == 0 {
		return "[]"
	}
	b, err := json.Marshal(items)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func decodeStringSliceJSON(raw string) []string {
	if raw == "" {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return []string{}
	}
	return out
}

func languageToCode(lang pb.Language) string {
	if lang == pb.Language_German {
		return "de"
	}
	return "en"
}

func codeToLanguage(code string) pb.Language {
	if code == "de" {
		return pb.Language_German
	}
	return pb.Language_English
}

func (s *CampusServer) GetPushPreferences(ctx context.Context, req *pb.GetPushPreferencesRequest) (*pb.PushPreferences, error) {
	s.recordDeviceIfPresent(ctx)
	if req.GetDeviceId() == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id is required.")
	}
	var dev model.PushDevice
	if err := s.db.WithContext(ctx).Select("device_id").First(&dev, "device_id = ?", req.DeviceId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Device not registered.")
		}
		log.WithError(err).Error("while checking push device")
		return nil, status.Error(codes.Internal, "Could not load push preferences.")
	}

	var pref model.PushPreference
	err := s.db.WithContext(ctx).First(&pref, "device_id = ?", req.DeviceId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &pb.PushPreferences{DeviceId: req.DeviceId, Language: pb.Language_English}, nil
	}
	if err != nil {
		log.WithError(err).Error("while loading push preferences")
		return nil, status.Error(codes.Internal, "Could not load push preferences.")
	}
	return pushPrefToProto(&pref), nil
}

func (s *CampusServer) UpdatePushPreferences(ctx context.Context, req *pb.UpdatePushPreferencesRequest) (*pb.PushPreferences, error) {
	s.recordDeviceIfPresent(ctx)
	if req.GetDeviceId() == "" {
		return nil, status.Error(codes.InvalidArgument, "device_id is required.")
	}
	canteens := normalizeStringList(req.CanteenIds, model.CanteenEatAPIID)
	favorites := normalizeStringList(req.FavoriteDishes, nil)
	dietary := normalizeStringList(req.DietaryLabels, strings.ToUpper)
	if len(canteens) > maxCanteenPrefs {
		return nil, status.Errorf(codes.InvalidArgument, "At most %d canteens can be monitored.", maxCanteenPrefs)
	}
	if len(favorites) > maxFavoriteDishKeywords {
		return nil, status.Errorf(codes.InvalidArgument, "At most %d favorite dishes can be configured.", maxFavoriteDishKeywords)
	}
	if len(dietary) > maxDietaryLabels {
		return nil, status.Errorf(codes.InvalidArgument, "At most %d dietary labels can be configured.", maxDietaryLabels)
	}

	tx := s.db.WithContext(ctx)
	var dev model.PushDevice
	if err := tx.Select("device_id").First(&dev, "device_id = ?", req.DeviceId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Device must be registered via CreateDevice first.")
		}
		log.WithError(err).Error("while verifying push device before updating preferences")
		return nil, status.Error(codes.Internal, "Could not update push preferences.")
	}

	pref := model.PushPreference{
		DeviceID:             req.DeviceId,
		CanteenIDs:           encodeStringSliceJSON(canteens),
		FavoriteDishes:       encodeStringSliceJSON(favorites),
		DietaryLabels:        encodeStringSliceJSON(dietary),
		NotifyFavoriteDishes: req.NotifyFavoriteDishes,
		NotifyITStatus:       req.NotifyItStatus,
		NotifyTUFilm:         req.NotifyTuFilm,
		Language:             languageToCode(req.Language),
		UpdatedAt:            time.Now(),
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "device_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"canteen_ids", "favorite_dishes", "dietary_labels",
			"notify_favorite_dishes", "notify_it_status", "notify_tu_film",
			"language", "updated_at",
		}),
	}).Create(&pref).Error; err != nil {
		log.WithError(err).Error("while saving push preferences")
		return nil, status.Error(codes.Internal, "Could not update push preferences.")
	}
	return pushPrefToProto(&pref), nil
}

func pushPrefToProto(p *model.PushPreference) *pb.PushPreferences {
	return &pb.PushPreferences{
		DeviceId:             p.DeviceID,
		CanteenIds:           decodeStringSliceJSON(p.CanteenIDs),
		FavoriteDishes:       decodeStringSliceJSON(p.FavoriteDishes),
		DietaryLabels:        decodeStringSliceJSON(p.DietaryLabels),
		NotifyFavoriteDishes: p.NotifyFavoriteDishes,
		NotifyItStatus:       p.NotifyITStatus,
		NotifyTuFilm:         p.NotifyTUFilm,
		Language:             codeToLanguage(p.Language),
	}
}
