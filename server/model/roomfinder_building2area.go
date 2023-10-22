package model

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

// RoomfinderBuilding2area struct is a row record of the roomfinder_building2area table in the tca database
type RoomfinderBuilding2area struct {
	BuildingNr string `gorm:"primary_key;column:building_nr;type:varchar(8);" json:"building_nr"`
	AreaID     int32  `gorm:"column:area_id;type:int;" json:"area_id"`
	Campus     string `gorm:"column:campus;type:char;size:1;" json:"campus"`
	Name       string `gorm:"column:name;type:varchar(32);" json:"name"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderBuilding2area) TableName() string {
	return "roomfinder_building2area"
}
