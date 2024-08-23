package model

import (
	"time"
)

type DishRating struct {
	DishRating int64     `gorm:"primary_key;AUTO_INCREMENT;column:dishRating;type:int;not null;" json:"dishRating"`
	Points     int32     `gorm:"column:points;type:int;not null;" json:"points"`
	DishID     int64     `gorm:"column:dishID;type:int;not null;index:dish_rating_dish_dish_fk" json:"dishID"`
	Dish       Dish      `gorm:"foreignKey:dishID;references:dish;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Comment    string    `gorm:"column:comment;type:text;" json:"comment"`
	Timestamp  time.Time `gorm:"column:timestamp;type:timestamp;not null;default:current_timestamp();OnUpdate:current_timestamp();" json:"timestamp"`
	Image      string    `gorm:"column:image;type:text;" json:"image"`
}
