package cafeteria_rating_models

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

// TagRating struct is a row record of the either the meal_tag_rating-table or the cafeteria_rating_tags-table in the database
type TagRating struct {
	Id           int32  `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id" :"id"`
	ParentRating int32  `gorm:"foreignKey:cafeteriaRatingID;column:parentRating;type:int;" json:"parentRating" :"id"`
	Rating       int32  `gorm:"column:rating;type:int;" json:"rating" :"rating"`
	Tagname      string `gorm:"column:tagname;type:varchar;size:32" json:"tagname" :"tagname"`
}
