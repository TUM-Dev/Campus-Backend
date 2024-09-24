package model

// DishNameTagStatistic is a view for statistics of DishName ratings
type DishNameTagStatistic struct {
	CafeteriaID int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"cafeteriaID"`
	TagID       int64   `gorm:"column:tagID;foreignKey:DishNameTagOption;type:int;not null;" json:"tagID"`
	Average     float64 `gorm:"column:average;type:double;not null;" json:"average" `
	Min         int8    `gorm:"column:min;type:int;not null;" json:"min"`
	Max         int8    `gorm:"column:max;type:int;not null;" json:"max"`
	Std         float64 `gorm:"column:std;type:double;not null;" json:"std"`
}
