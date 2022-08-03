package model

// DishNameTagAverage stores all precomputed values for the DishName ratings
type DishNameTagAverage struct {
	DishNameTagAverage int32   `gorm:"primary_key;AUTO_INCREMENT;column:dishNameTagAverage;type:int;not null;" json:"dishNameTagAverage" `
	CafeteriaID        int32   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"cafeteriaID"`
	TagID              int32   `gorm:"column:tagID;foreignKey:DishNameTagOption;type:int;not null;" json:"tagID"`
	Average            float32 `gorm:"column:average;type:float;not null;" json:"average" `
	Min                int8    `gorm:"column:min;type:int;not null;" json:"min"`
	Max                int8    `gorm:"column:max;type:int;not null;" json:"max"`
	Std                float32 `gorm:"column:std;type:float;not null;" json:"std"`
}

// TableName sets the insert table name for this struct type
func (n *DishNameTagAverage) TableName() string {
	return "dish_name_tag_average"
}
