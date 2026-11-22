package model

import "time"

// PushPreference stores per-device opt-in notification settings for Mensa dishes, IT outages, and TU-Film.
type PushPreference struct {
	DeviceID             string    `gorm:"primaryKey;column:device_id;type:varchar(191);not null;" json:"device_id"`
	CanteenIDs           string    `gorm:"column:canteen_ids;type:text;not null;" json:"canteen_ids"`
	FavoriteDishes       string    `gorm:"column:favorite_dishes;type:text;not null;" json:"favorite_dishes"`
	DietaryLabels        string    `gorm:"column:dietary_labels;type:text;not null;" json:"dietary_labels"`
	NotifyFavoriteDishes bool      `gorm:"column:notify_favorite_dishes;type:tinyint(1);not null;default:0;" json:"notify_favorite_dishes"`
	NotifyITStatus       bool      `gorm:"column:notify_it_status;type:tinyint(1);not null;default:0;" json:"notify_it_status"`
	NotifyTUFilm         bool      `gorm:"column:notify_tu_film;type:tinyint(1);not null;default:0;" json:"notify_tu_film"`
	Language             string    `gorm:"column:language;type:enum('de','en');not null;default:'en';" json:"language"`
	UpdatedAt            time.Time `gorm:"column:updated_at;type:timestamp;not null;default:current_timestamp() ON UPDATE current_timestamp();" json:"updated_at"`
}
