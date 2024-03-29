package model

import (
	"time"
)

// CafeteriaRating stores all Available cafeterias in the format of the eat-api
type CafeteriaRating struct {
	CafeteriaRating int64     `gorm:"primary_key;AUTO_INCREMENT;column:cafeteriaRating;type:int;not null;" json:"canteenrating"`
	Points          int32     `gorm:"column:points;type:int;not null;" json:"points"`
	Comment         string    `gorm:"column:comment;type:text;" json:"comment" `
	CafeteriaID     int64     `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"canteenID"`
	Timestamp       time.Time `gorm:"column:timestamp;type:timestamp;not null;" json:"timestamp" `
	Image           string    `gorm:"column:image;type:text;" json:"image"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRating) TableName() string {
	return "cafeteria_rating"
}
