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

// News struct is a row record of the news table in the tca database
type MealRating struct {
	Id        int32     `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id"`
	Rating    int32     `gorm:"column:rating;type:int;" json:"rating"`
	Meal      string    `gorm:"column:meal;type:varchar;size:128;" json:"meal"`
	Comment   string    `gorm:"column:comment;type:varchar;size:256;" json:"comment"`
	Cafeteria string    `gorm:"column:canteen;type:mediumtext;" json:"canteen"`
	Timestamp time.Time `gorm:"column:timestamp;type:timestamp;default:CURRENT_TIMESTAMP;" json:"timestamp"`
	//	Image     null.String `gorm:"column:image;type:text;size:65535;" json:"image" :"image"`
}

// TableName sets the insert table name for this struct type
func (n *MealRating) TableName() string {
	return "meal_rating"
}
