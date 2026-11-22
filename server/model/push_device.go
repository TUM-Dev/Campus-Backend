package model

import "time"

// PushDevice stores a client's Firebase Cloud Messaging token for push notifications.
type PushDevice struct {
	DeviceID   string          `gorm:"primaryKey;column:device_id;type:varchar(191);not null;" json:"device_id"`
	DeviceType string          `gorm:"column:device_type;type:enum('ios','android','windows');not null;" json:"device_type"`
	FCMToken   string          `gorm:"column:fcm_token;type:varchar(512);not null;" json:"fcm_token"`
	UpdatedAt  time.Time       `gorm:"column:updated_at;type:timestamp;not null;default:current_timestamp() ON UPDATE current_timestamp();" json:"updated_at"`
	Preference *PushPreference `gorm:"foreignKey:DeviceID;references:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
