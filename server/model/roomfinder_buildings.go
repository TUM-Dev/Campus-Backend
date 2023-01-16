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

// RoomfinderBuildings struct is a row record of the roomfinder_buildings table in the tca database
type RoomfinderBuildings struct {
	//[ 0] building_nr                                    varchar(8)           null: false  primary: true   isArray: false  auto: false  col: varchar         len: 8       default: []
	BuildingNr string `gorm:"primary_key;column:building_nr;type:varchar(8);" json:"building_nr"`
	//[ 1] utm_zone                                       varchar(4)           null: true   primary: false  isArray: false  auto: false  col: varchar         len: 4       default: []
	UtmZone null.String `gorm:"column:utm_zone;type:varchar(4);" json:"utm_zone"`
	//[ 2] utm_easting                                    varchar(32)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	UtmEasting null.String `gorm:"column:utm_easting;type:varchar(32);" json:"utm_easting"`
	//[ 3] utm_northing                                   varchar(32)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	UtmNorthing null.String `gorm:"column:utm_northing;type:varchar(32);" json:"utm_northing"`
	//[ 4] default_map_id                                 int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	DefaultMapID null.Int `gorm:"column:default_map_id;type:int;" json:"default_map_id"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderBuildings) TableName() string {
	return "roomfinder_buildings"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *RoomfinderBuildings) BeforeSave() error {
	return nil
}
