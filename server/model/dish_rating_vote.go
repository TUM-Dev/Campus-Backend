package model

import "time"

// DishRatingVote records a "helpful" upvote cast by a device on a dish rating.
type DishRatingVote struct {
	ID         int64      `gorm:"primary_key;autoIncrement;column:id;type:bigint;not null;" json:"id"`
	RatingID   int64      `gorm:"column:rating_id;type:bigint;not null;uniqueIndex:idx_dish_rating_votes_unique,priority:1;index:idx_dish_rating_votes_rating" json:"rating_id"`
	Rating     DishRating `gorm:"foreignKey:RatingID;references:DishRating;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DeviceHash string     `gorm:"column:device_hash;type:varchar(64);not null;uniqueIndex:idx_dish_rating_votes_unique,priority:2" json:"device_hash"`
	CreatedAt  time.Time  `gorm:"column:created_at;type:timestamp;not null;default:current_timestamp();" json:"created_at"`
}
