package model

import (
	"time"
)

type DishRating struct {
	DishRating  int64     `gorm:"primary_key;AUTO_INCREMENT;column:dishRating;type:int;not null;" json:"dishRating"`
	Points      int32     `gorm:"column:points;type:int;not null;" json:"points"`
	CafeteriaID int64     `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"cafeteriaID"`
	DishID      int64     `gorm:"column:dishID;foreignKey:dish;type:int;not null;" json:"dishID"`
	Comment     string    `gorm:"column:comment;type:text;" json:"comment"`
	Timestamp   time.Time `gorm:"column:timestamp;type:timestamp;not null;" json:"timestamp"`
	Image       string    `gorm:"column:image;type:text;" json:"image"`
}
