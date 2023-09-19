package model

// CafeteriaRatingAverage stores all precomputed values for the cafeteria ratings
type CafeteriaRatingAverage struct {
	CafeteriaRatingAverage int64   `gorm:"primary_key;AUTO_INCREMENT;column:cafeteriaRatingAverage;type:int;not null;" json:"canteenRatingAverage" `
	CafeteriaID            int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"canteenID"`
	Average                float64 `gorm:"column:average;type:float;not null;" json:"average" `
	Min                    int32   `gorm:"column:min;type:int;not null;" json:"min"`
	Max                    int32   `gorm:"column:max;type:int;not null;" json:"max"`
	Std                    float64 `gorm:"column:std;type:float;not null;" json:"std"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRatingAverage) TableName() string {
	return "cafeteria_rating_average"
}
