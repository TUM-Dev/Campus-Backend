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

// Cafeteria stores all Available cafeterias in the format of the eat-api
type CafeteriaRating struct {
	CafeteriaRating int32     `gorm:"primary_key;AUTO_INCREMENT;column:cafeteriaRating;type:int;" json:"cafeteriarating"`
	Points          int32     `gorm:"column:points;type:int;" json:"points"`
	Comment         string    `gorm:"column:comment;type:text;" json:"comment" `
	CafeteriaID     int32     `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	Timestamp       time.Time `gorm:"column:timestamp;type:timestamp;" json:"timestamp" `
	Image           string    `gorm:"column:image;type:text;" json:"image"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRating) TableName() string {
	return "cafeteria_rating"
}
