package cafeteria_rating_models

import (
	"database/sql"
	"github.com/guregu/null"
)

var (
	_ = sql.LevelDefault
	_ = null.Bool{}
)

// News struct is a row record of the mensa table in the tca database
type CafeteriaRatingTagsAverage struct {
	Id          int32   `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id" `
	CafeteriaID int32   `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	TagID       int32   `gorm:"column:tagID;foreignKey:tagID;type:int;" json:"tagID"`
	Average     float32 `gorm:"column:average;type:float;" json:"average" `
	Min         int8    `gorm:"column:min;type:int;" json:"min"`
	Max         int8    `gorm:"column:max;type:int;" json:"max"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRatingTagsAverage) TableName() string {
	return "cafeteria_rating_tags_results"
}
