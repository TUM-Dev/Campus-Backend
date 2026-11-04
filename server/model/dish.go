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
	Dish        int64  `gorm:"primary_key;autoIncrement;column:dish;type:int;not null;" json:"dish"`
	Name        string `gorm:"column:name;type:text;not null;uniqueIndex:dish_name_cafeteriaID_uindex,expression:name(255)" json:"name" `
	Type        string `gorm:"column:type;type:text;not null;" json:"type" `
	CafeteriaID int64  `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;uniqueIndex:dish_name_cafeteriaID_uindex" json:"cafeteriaID"`
}

// TableName sets the insert table name for this struct type
func (n *Dish) TableName() string {
	return "dishes"
}
