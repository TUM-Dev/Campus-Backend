package model

import (
	"time"
)

// CanteenRating stores all Available cafeterias in the format of the eat-api
type CanteenRating struct {
	CafeteriaRating int64     `gorm:"primary_key;autoIncrement;column:cafeteriaRating;type:int;not null;" json:"canteenrating"`
	Points          int32     `gorm:"column:points;type:int;not null;" json:"points"`
	Comment         string    `gorm:"column:comment;type:text;" json:"comment" `
	CafeteriaID     int64     `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"canteenID"`
	Timestamp       time.Time `gorm:"column:timestamp;type:timestamp;not null;default:current_timestamp();OnUpdate:current_timestamp();" json:"timestamp" `
	Image           string    `gorm:"column:image;type:text;" json:"image"`
}
