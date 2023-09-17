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

// RoomfinderMaps struct is a row record of the roomfinder_maps table in the tca database
type RoomfinderMaps struct {
	MapID       int32  `gorm:"primary_key;column:map_id;type:int;" json:"map_id"`
	Description string `gorm:"column:description;type:varchar(64);" json:"description"`
	Scale       int32  `gorm:"column:scale;type:int;" json:"scale"`
	Width       int32  `gorm:"column:width;type:int;" json:"width"`
	Height      int32  `gorm:"column:height;type:int;" json:"height"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderMaps) TableName() string {
	return "roomfinder_maps"
}
