package model

// CanteenRatingStatistic is a view for statistics of canteen ratings
type CanteenRatingStatistic struct {
	CanteenID int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;"`
	Average   float64 `gorm:"column:average;type:double;not null;"`
	Min       int32   `gorm:"column:min;type:int;not null;"`
	Max       int32   `gorm:"column:max;type:int;not null;"`
	Std       float64 `gorm:"column:std;type:double;not null;"`
}

// TableName pins the legacy view name, which predates the Canteen rename.
func (CanteenRatingStatistic) TableName() string {
	return "cafeteria_rating_statistics"
}
