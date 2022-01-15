package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	"github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

// RoomfinderBuildings2maps struct is a row record of the roomfinder_buildings2maps table in the tca database
type RoomfinderBuildings2maps struct {
	//[ 0] building_nr                                    varchar(8)           null: false  primary: true   isArray: false  auto: false  col: varchar         len: 8       default: []
	BuildingNr string `gorm:"primary_key;column:building_nr;type:varchar;size:8;" json:"building_nr"`
	//[ 1] map_id                                         int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	MapID int32 `gorm:"primary_key;column:map_id;type:int;" json:"map_id"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderBuildings2maps) TableName() string {
	return "roomfinder_buildings2maps"
}
