package model

import (
	"time"
)

// CanteenHeadCount stores all available people counts for available canteens. The CanteenId represents the same ID, as for the canteen inside the eat-api.
type CanteenHeadCount struct {
	CanteenId string    `gorm:"primary_key;column:canteen_id;type:varchar(64);not null;" json:"canteen_id"`
	Count     uint32    `gorm:"column:count;type:int;not null;" json:"count"`
	MaxCount  uint32    `gorm:"column:max_count;type:int;not null;" json:"max_count"`
	Percent   float32   `gorm:"column:percent;type:float;not null;" json:"percent"`
	Timestamp time.Time `gorm:"column:timestamp;type:timestamp;not null;default:current_timestamp();OnUpdate:current_timestamp();" json:"timestamp" `
}
