package model

import (
	"database/sql"

	"github.com/guregu/null"
)

var (
	_ = sql.LevelDefault
	_ = null.Bool{}
)

// Dish represents one dish fin a specific cafeteria
type Dish struct {
	Dish        int32  `gorm:"primary_key;AUTO_INCREMENT;column:dish;type:int;not null;" json:"dish"`
	Name        string `gorm:"column:name;type:text;not null;" json:"name" `
	Type        string `gorm:"column:type;type:text;not null;" json:"type" `
	CafeteriaID int32  `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"cafeteriaID"`
}

// TableName sets the insert table name for this struct type
func (n *Dish) TableName() string {
	return "dish"
}
