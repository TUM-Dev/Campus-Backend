package model

// DishRatingStatistic stores all precomputed values for the cafeteria ratings
type DishRatingStatistic struct {
	CafeteriaID int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;"`
	DishID      int64   `gorm:"column:dishID;foreignKey:dish;type:int;not null;"`
	Average     float64 `gorm:"column:average;type:float;not null;"`
	Min         int32   `gorm:"column:min;type:int;not null;"`
	Max         int32   `gorm:"column:max;type:int;not null;"`
	Std         float64 `gorm:"column:std;type:float;not null;"`
}
