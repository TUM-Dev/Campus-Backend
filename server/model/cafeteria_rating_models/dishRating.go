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

type DishRating struct {
	DishRating  int32     `gorm:"primary_key;AUTO_INCREMENT;column:dishRating;type:int;" json:"dishRating"`
	Points      int32     `gorm:"column:points;type:int;" json:"points"`
	CafeteriaID int32     `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;" json:"cafeteriaID"`
	DishID      int32     `gorm:"column:dishID;foreignKey:dish;type:int;" json:"dishID"`
	Comment     string    `gorm:"column:comment;type:text" json:"comment"`
	Timestamp   time.Time `gorm:"column:timestamp;type:timestamp;" json:"timestamp"`
	Image       string    `gorm:"column:image;type:text;" json:"image"`
}

// TableName sets the insert table name for this struct type
func (n *DishRating) TableName() string {
	return "dish_rating"
}
