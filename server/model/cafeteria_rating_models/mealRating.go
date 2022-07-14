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

type MealRating struct {
	MealRating  int32     `gorm:"primary_key;AUTO_INCREMENT;column:mealRating;type:int;" json:"mealRating"`
	Points      int32     `gorm:"column:points;type:int;" json:"points"`
	CafeteriaID int32     `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	MealID      int32     `gorm:"column:mealID;foreignKey:id;type:int;" json:"mealID"`
	Comment     string    `gorm:"column:comment;type:varchar;size:256;" json:"comment"`
	Timestamp   time.Time `gorm:"column:timestamp;type:timestamp;" json:"timestamp"`
	Image       string    `gorm:"column:image;type:text;" json:"image"`
}

// TableName sets the insert table name for this struct type
func (n *MealRating) TableName() string {
	return "meal_rating"
}
