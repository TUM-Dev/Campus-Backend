package model

// CanteenRatingTagsStatistic is a view for statistics of canteen rating tags
type CanteenRatingTagsStatistic struct {
	CanteenID int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"canteenID"`
	TagID     int64   `gorm:"column:tagID;foreignKey:cafeteriaRatingTagOption;type:int;not null;" json:"tagID"`
	Average   float64 `gorm:"column:average;type:double;not null;" json:"average"`
	Min       int8    `gorm:"column:min;type:int;not null;" json:"min"`
	Max       int8    `gorm:"column:max;type:int;not null;" json:"max"`
	Std       float64 `gorm:"column:std;type:double;not null;" json:"std"`
}

// TableName pins the legacy view name, which predates the Canteen rename.
func (CanteenRatingTagsStatistic) TableName() string {
	return "cafeteria_rating_tag_statistics"
}
