package mensa_rating_models

import (
	"database/sql"
	"github.com/guregu/null"
	"time"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

// News struct is a row record of the news table in the tca database
type TagRating struct {
	Id           int32  `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id" :"id"`
	ParentRating int32  `gorm:"foreignKey:cafeteriaRatingID;column:parentRating;type:int;" json:"id" :"id"`
	Rating       int32  `gorm:"column:rating;type:int;" json:"rating" :"rating"`
	Tagname      string `gorm:"column:tagname;type:varchar;size:32" json:"tagname" :"tagname"`
}
