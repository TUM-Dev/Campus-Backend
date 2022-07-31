package model

// DishRatingsAverage stores all precomputed values for the cafeteria ratings
type DishRatingAverage struct {
	DishRatingAverage int32   `gorm:"primary_key;AUTO_INCREMENT;column:dishRatingAverage;type:int;" json:"dishRatingAverage" `
	CafeteriaID       int32   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;" json:"cafeteriaID"`
	DishID            int32   `gorm:"column:dishID;foreignKey:dish;type:int;" json:"dishID"`
	Average           float32 `gorm:"column:average;type:float;" json:"average" `
	Min               int8    `gorm:"column:min;type:int;" json:"min"`
	Max               int8    `gorm:"column:max;type:int;" json:"max"`
	Std               float32 `gorm:"column:std;type:float;" json:"std"`
}

// TableName sets the insert table name for this struct type
func (n *DishRatingAverage) TableName() string {
	return "dish_rating_average"
}
