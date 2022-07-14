package cafeteria_rating_models

// MealRatingsAverage stores all precomputed values for the cafeteria ratings
type MealRatingAverage struct {
	MealRatingAverage int32   `gorm:"primary_key;AUTO_INCREMENT;column:mealRatingAverage;type:int;" json:"mealRatingAverage" `
	CafeteriaID       int32   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;" json:"cafeteriaID"`
	MealID            int32   `gorm:"column:mealID;foreignKey:meal;type:int;" json:"mealID"`
	Average           float32 `gorm:"column:average;type:float;" json:"average" `
	Min               int8    `gorm:"column:min;type:int;" json:"min"`
	Max               int8    `gorm:"column:max;type:int;" json:"max"`
	Std               float32 `gorm:"column:std;type:float;" json:"std"`
}

// TableName sets the insert table name for this struct type
func (n *MealRatingAverage) TableName() string {
	return "meal_rating_average"
}
