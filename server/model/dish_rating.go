package model

import (
	"time"
)

type DishRating struct {
	DishRating int64     `gorm:"primary_key;AUTO_INCREMENT;column:dishRating;type:int;not null;" json:"dishRating"`
	Points     int32     `gorm:"column:points;type:int;not null;" json:"points"`
	DishID     int64     `gorm:"column:dishID;type:int;not null;" json:"dishID"`
	Dish       Dish      `gorm:"foreignKey:dishID;references:dish;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comment    string    `gorm:"column:comment;type:text;" json:"comment"`
	Timestamp  time.Time `gorm:"column:timestamp;type:timestamp;not null;default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP;" json:"timestamp"`
	Image      string    `gorm:"column:image;type:text;" json:"image"`
}

// TableName sets the insert table name for this struct type
func (n *DishRating) TableName() string {
	return "dish_rating"
}
