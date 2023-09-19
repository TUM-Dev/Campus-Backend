package model

// DishRatingAverage stores all precomputed values for the cafeteria ratings
type DishRatingAverage struct {
	DishRatingAverage int64   `gorm:"primary_key;AUTO_INCREMENT;column:dishRatingAverage;type:int;not null;" json:"dishRatingAverage" `
	CafeteriaID       int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"cafeteriaID"`
	DishID            int64   `gorm:"column:dishID;foreignKey:dish;type:int;not null;" json:"dishID"`
	Average           float32 `gorm:"column:average;type:float;not null;" json:"average" `
	Min               int8    `gorm:"column:min;type:int;not null;" json:"min"`
	Max               int8    `gorm:"column:max;type:int;not null;" json:"max"`
	Std               float32 `gorm:"column:std;type:float;not null;" json:"std"`
}

// TableName sets the insert table name for this struct type
func (n *DishRatingAverage) TableName() string {
	return "dish_rating_average"
}
