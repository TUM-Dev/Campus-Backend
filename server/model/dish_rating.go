package model

import (
	"time"

	"github.com/guregu/null"
)

type DishRating struct {
	DishRating int64       `gorm:"primary_key;autoIncrement;column:dishRating;type:int;not null;" json:"dishRating"`
	Points     int32       `gorm:"column:points;type:int;not null;" json:"points"`
	DishID     int64       `gorm:"column:dishID;type:int;not null;index:dish_rating_dish_dish_fk;index:idx_dish_ratings_device_dish_ts,priority:2" json:"dishID"`
	Dish       Dish        `gorm:"foreignKey:dishID;references:dish;constraint:-;"`
	Comment    string      `gorm:"column:comment;type:text;" json:"comment"`
	Timestamp  time.Time   `gorm:"column:timestamp;type:timestamp;not null;default:current_timestamp() ON UPDATE current_timestamp();index:idx_dish_ratings_device_dish_ts,priority:3" json:"timestamp"`
	Image      string      `gorm:"column:image;type:text;" json:"image"`
	DeviceHash null.String `gorm:"column:device_hash;type:varchar(64);index:idx_dish_ratings_device_dish_ts,priority:1" json:"device_hash"`
	Status     string      `gorm:"column:status;type:enum('visible','hidden');not null;default:'visible';index:idx_dish_ratings_status" json:"status"`
	UpdatedAt  null.Time   `gorm:"column:updated_at;type:timestamp;autoUpdateTime:false;" json:"updated_at"`
}
