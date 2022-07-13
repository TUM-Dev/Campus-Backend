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
type CafeteriaRating struct {
	Id          int32     `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id"`
	Rating      int32     `gorm:"column:rating;type:int;" json:"rating"`
	Comment     string    `gorm:"column:comment;type:varchar;size:256;" json:"comment" `
	CafeteriaID int32     `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	Timestamp   time.Time `gorm:"column:timestamp;type:timestamp;default:CURRENT_TIMESTAMP;" json:"timestamp" `
	Image       string    `gorm:"column:image;type:mediumtext;" json:"image"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRating) TableName() string {
	return "cafeteria_rating"
}
