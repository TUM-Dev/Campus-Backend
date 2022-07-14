package cafeteria_rating_models

import (
	"database/sql"
	"github.com/guregu/null"
)

var (
	_ = sql.LevelDefault
	_ = null.Bool{}
)

// Meal represents one meal fin a specific cafeteria
type Meal struct {
	Meal        int32  `gorm:"primary_key;AUTO_INCREMENT;column:meal;type:int;" json:"meal"`
	Name        string `gorm:"column:name;type:text;" json:"name" `
	Type        string `gorm:"column:type;type:text;" json:"type" `
	CafeteriaID int32  `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;" json:"cafeteriaID"`
}

// TableName sets the insert table name for this struct type
func (n *Meal) TableName() string {
	return "meal"
}
