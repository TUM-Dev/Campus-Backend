package model

// CafeteriaRatingStatistic is a view for statistics of cafeteria ratings
type CafeteriaRatingStatistic struct {
	CafeteriaID int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;"`
	Average     float64 `gorm:"column:average;type:double;not null;"`
	Min         int32   `gorm:"column:min;type:int;not null;"`
	Max         int32   `gorm:"column:max;type:int;not null;"`
	Std         float64 `gorm:"column:std;type:double;not null;"`
}
