package model

import (
	"database/sql"
	"github.com/guregu/null"
)

var (
	_ = sql.LevelDefault
	_ = null.Bool{}
)

// News struct is a row record of the mensa table in the tca database
type Dish struct {
	id      int32  `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id"`
	Name    string `gorm:"column:name;type:varchar;size:150;;" json:"name" `
	Type    string `gorm:"column:type;type:varchar;size:20;" json:"type" `
	Canteen string `gorm:"column:canteen;type:mediumtext;" json:"canteen" `
}

// TableName sets the insert table name for this struct type
func (n *Dish) TableName() string {
	return "dish"
}
