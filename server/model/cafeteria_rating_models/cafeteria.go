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
type Cafeteria struct {
	Id        int32   `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id" `
	Name      string  `gorm:"column:name;type:mediumtext;" json:"name" `
	Address   string  `gorm:"column:address;type:mediumtext" json:"address" `
	Latitude  float32 `gorm:"column:latitude;type:float;" json:"latitude" `
	Longitude float32 `gorm:"column:longitude;type:float;" json:"longitude"`
}

// TableName sets the insert table name for this struct type
func (n *Cafeteria) TableName() string {
	return "cafeteria"
}
