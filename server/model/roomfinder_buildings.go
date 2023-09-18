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

// RoomfinderBuildings struct is a row record of the roomfinder_buildings table in the tca database
type RoomfinderBuildings struct {
	BuildingNr   string      `gorm:"primary_key;column:building_nr;type:varchar(8);" json:"building_nr"`
	UtmZone      null.String `gorm:"column:utm_zone;type:varchar(4);" json:"utm_zone"`
	UtmEasting   null.String `gorm:"column:utm_easting;type:varchar(32);" json:"utm_easting"`
	UtmNorthing  null.String `gorm:"column:utm_northing;type:varchar(32);" json:"utm_northing"`
	DefaultMapID null.Int    `gorm:"column:default_map_id;type:int;" json:"default_map_id"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderBuildings) TableName() string {
	return "roomfinder_buildings"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *RoomfinderBuildings) BeforeSave() error {
	return nil
}
