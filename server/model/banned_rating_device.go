package model

import "time"

// BannedRatingDevice blocks an abusive hashed device ID from submitting ratings, votes, or reports.
type BannedRatingDevice struct {
	DeviceHash string    `gorm:"primaryKey;column:device_hash;type:varchar(64);not null;" json:"device_hash"`
	Reason     string    `gorm:"column:reason;type:varchar(256);not null;default:'';" json:"reason"`
	BannedAt   time.Time `gorm:"column:banned_at;type:timestamp;not null;default:current_timestamp();" json:"banned_at"`
}
