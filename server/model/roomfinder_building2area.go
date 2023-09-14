package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

// RoomfinderBuilding2area struct is a row record of the roomfinder_building2area table in the tca database
type RoomfinderBuilding2area struct {
	//[ 0] area_id                                        int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	AreaID int32 `gorm:"column:area_id;type:int;" json:"area_id"`
	//[ 1] building_nr                                    varchar(8)           null: false  primary: true   isArray: false  auto: false  col: varchar         len: 8       default: []
	BuildingNr string `gorm:"primary_key;column:building_nr;type:varchar(8);" json:"building_nr"`
	//[ 2] campus                                         char(1)              null: false  primary: false  isArray: false  auto: false  col: char            len: 1       default: []
	Campus string `gorm:"column:campus;type:char;size:1;" json:"campus"`
	//[ 3] name                                           varchar(32)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	Name string `gorm:"column:name;type:varchar(32);" json:"name"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderBuilding2area) TableName() string {
	return "roomfinder_building2area"
}
