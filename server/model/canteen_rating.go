package model

import (
	"time"

	"github.com/guregu/null"
)

// CanteenRating is a single user rating of a canteen
type CanteenRating struct {
	CanteenRating int64       `gorm:"primary_key;autoIncrement;column:cafeteriaRating;type:int;not null;" json:"canteenrating"`
	Points        int32       `gorm:"column:points;type:int;not null;" json:"points"`
	Comment       string      `gorm:"column:comment;type:text;" json:"comment" `
	CanteenID     int64       `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;index:cafeteria_rating_cafeteria_cafeteria_fk;" json:"canteenID"`
	Timestamp     time.Time   `gorm:"column:timestamp;type:timestamp;not null;default:current_timestamp() ON UPDATE current_timestamp();" json:"timestamp" `
	Image         string      `gorm:"column:image;type:text;" json:"image"`
	DeviceHash    null.String `gorm:"column:device_hash;type:varchar(64);index:idx_canteen_ratings_device_ts" json:"device_hash"`
	Status        string      `gorm:"column:status;type:enum('visible','hidden');not null;default:'visible';index:idx_canteen_ratings_status" json:"status"`
	UpdatedAt     null.Time   `gorm:"column:updated_at;type:timestamp;autoUpdateTime:false;" json:"updated_at"`
}
