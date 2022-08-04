package model

// CafeteriaRatingTagsAverage stores all precomputed values for the cafeteria ratings
type CafeteriaRatingTagsAverage struct {
	CafeteriaRatingTagsAverage int32   `gorm:"primary_key;AUTO_INCREMENT;column:cafeteriaRatingTagsAverage;type:int;not null;" json:"cafeteriaRatingTagsAverage" `
	CafeteriaID                int32   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"cafeteriaID"`
	TagID                      int32   `gorm:"column:tagID;foreignKey:cafeteriaRatingTagOption;type:int;not null;" json:"tagID"`
	Average                    float32 `gorm:"column:average;type:float;not null;" json:"average" `
	Min                        int8    `gorm:"column:min;type:int;not null;" json:"min"`
	Max                        int8    `gorm:"column:max;type:int;not null;" json:"max"`
	Std                        float32 `gorm:"column:std;type:float;not null;" json:"std"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRatingTagsAverage) TableName() string {
	return "cafeteria_rating_tag_average"
}
