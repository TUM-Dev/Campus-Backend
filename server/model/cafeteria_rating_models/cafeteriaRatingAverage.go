package cafeteria_rating_models

// CafeteriaRatingAverage stores all precomputed values for the cafeteria ratings
type CafeteriaRatingAverage struct {
	CafeteriaRatingAverage int32   `gorm:"primary_key;AUTO_INCREMENT;column:cafeteriaRatingAverage;type:int;" json:"cafeteriaRatingAverage" `
	CafeteriaID            int32   `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	Average                float32 `gorm:"column:average;type:float;" json:"average" `
	Min                    int8    `gorm:"column:min;type:int;" json:"min"`
	Max                    int8    `gorm:"column:max;type:int;" json:"max"`
	Std                    float32 `gorm:"column:std;type:float;" json:"std"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRatingAverage) TableName() string {
	return "cafeteria_rating_average"
}
