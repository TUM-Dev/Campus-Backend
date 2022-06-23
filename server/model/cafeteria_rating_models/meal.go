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
type Meal struct {
	Id          int32  `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id"`
	Name        string `gorm:"column:name;type:varchar;size:150;;" json:"name" `
	Type        string `gorm:"column:type;type:varchar;size:20;" json:"type" `
	CafeteriaID int32  `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
}

// TableName sets the insert table name for this struct type
func (n *Meal) TableName() string {
	return "dish"
}
