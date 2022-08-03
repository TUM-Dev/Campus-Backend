package model

import (
	"time"
)

type DishRating struct {
	DishRating  int32     `gorm:"primary_key;AUTO_INCREMENT;column:dishRating;type:int;not null;" json:"dishRating"`
	Points      int32     `gorm:"column:points;type:int;not null;" json:"points"`
	CafeteriaID int32     `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"cafeteriaID"`
	DishID      int32     `gorm:"column:dishID;foreignKey:dish;type:int;not null;" json:"dishID"`
	Comment     string    `gorm:"column:comment;type:text;" json:"comment"`
	Timestamp   time.Time `gorm:"column:timestamp;type:timestamp;not null;" json:"timestamp"`
	Image       string    `gorm:"column:image;type:text;" json:"image"`
}

// TableName sets the insert table name for this struct type
func (n *DishRating) TableName() string {
	return "dish_rating"
}
