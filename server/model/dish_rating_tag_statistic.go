package model

// DishRatingTagStatistic stores all precomputed values for the cafeteria ratings
type DishRatingTagStatistic struct {
	CafeteriaID int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"cafeteriaID"`
	TagID       int64   `gorm:"column:tagID;foreignKey:tagID;type:int;not null;" json:"tagID"`
	DishID      int64   `gorm:"column:dishID;foreignKey:dishID;type:int;not null;" json:"dishID"`
	Average     float32 `gorm:"column:average;type:float;not null;" json:"average" `
	Min         int8    `gorm:"column:min;type:int;not null;" json:"min"`
	Max         int8    `gorm:"column:max;type:int;not null;" json:"max"`
	Std         float32 `gorm:"column:std;type:float;not null;" json:"std"`
}
