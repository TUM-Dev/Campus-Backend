package cafeteria_rating_models

import (
	"database/sql"
	"github.com/guregu/null"
)

var (
	_ = sql.LevelDefault
	_ = null.Bool{}
)

// Cafeteria stores all Available cafeterias in the format of the eat-api
type Cafeteria struct {
	Cafeteria int32   `gorm:"primary_key;AUTO_INCREMENT;column:cafeteria;type:int;" json:"cafeteria" `
	Name      string  `gorm:"column:name;type:mediumtext;" json:"name" `
	Address   string  `gorm:"column:address;type:text" json:"address" `
	Latitude  float32 `gorm:"column:latitude;type:float;" json:"latitude" `
	Longitude float32 `gorm:"column:longitude;type:float;" json:"longitude"`
}

// TableName sets the insert table name for this struct type
func (n *Cafeteria) TableName() string {
	return "cafeteria"
}
