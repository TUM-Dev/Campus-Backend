package cafeteria_rating_models

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
	Dish        int32  `gorm:"primary_key;AUTO_INCREMENT;column:dish;type:int;" json:"dish"`
	Name        string `gorm:"column:name;type:text;" json:"name" `
	Type        string `gorm:"column:type;type:text;" json:"type" `
	CafeteriaID int32  `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;" json:"cafeteriaID"`
}

// TableName sets the insert table name for this struct type
func (n *Dish) TableName() string {
	return "dish"
}
