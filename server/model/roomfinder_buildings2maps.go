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

// RoomfinderBuildings2maps struct is a row record of the roomfinder_buildings2maps table in the tca database
type RoomfinderBuildings2maps struct {
	BuildingNr string `gorm:"primary_key;column:building_nr;type:varchar(8);" json:"building_nr"`
	MapID      int32  `gorm:"primary_key;column:map_id;type:int;" json:"map_id"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderBuildings2maps) TableName() string {
	return "roomfinder_buildings2maps"
}
