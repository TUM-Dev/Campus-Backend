package model

import "time"

const (
	RatingStatusVisible = "visible"
	RatingStatusHidden  = "hidden"
)

// DishRatingReport records a user report against a dish rating for moderation.
type DishRatingReport struct {
	ID         int64      `gorm:"primary_key;autoIncrement;column:id;type:bigint;not null;" json:"id"`
	RatingID   int64      `gorm:"column:rating_id;type:bigint;not null;uniqueIndex:idx_dish_rating_reports_rating_device" json:"rating_id"`
	Rating     DishRating `gorm:"foreignKey:RatingID;references:DishRating;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DeviceHash string     `gorm:"column:device_hash;type:varchar(64);not null;uniqueIndex:idx_dish_rating_reports_rating_device" json:"device_hash"`
	Reason     string     `gorm:"column:reason;type:varchar(256);not null;default:''" json:"reason"`
	CreatedAt  time.Time  `gorm:"column:created_at;type:timestamp;not null;default:current_timestamp();" json:"created_at"`
}
