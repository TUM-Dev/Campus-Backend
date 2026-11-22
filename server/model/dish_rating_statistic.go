package model

// DishRatingStatistic is a view for statistics of dish ratings
type DishRatingStatistic struct {
	CanteenID   int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;"`
	DishID      int64   `gorm:"column:dishID;foreignKey:dish;type:int;not null;"`
	Average     float64 `gorm:"column:average;type:double;not null;"`
	Min         int32   `gorm:"column:min;type:int;not null;"`
	Max         int32   `gorm:"column:max;type:int;not null;"`
	Std         float64 `gorm:"column:std;type:double;not null;"`
	RatingCount int64   `gorm:"column:rating_count;type:bigint;not null;"`
	Stars1      int64   `gorm:"column:stars_1;type:bigint;not null;"`
	Stars2      int64   `gorm:"column:stars_2;type:bigint;not null;"`
	Stars3      int64   `gorm:"column:stars_3;type:bigint;not null;"`
	Stars4      int64   `gorm:"column:stars_4;type:bigint;not null;"`
	Stars5      int64   `gorm:"column:stars_5;type:bigint;not null;"`
}
