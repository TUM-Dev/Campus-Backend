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

// RoomfinderMaps struct is a row record of the roomfinder_maps table in the tca database
type RoomfinderMaps struct {
	//[ 0] map_id                                         int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	MapID int32 `gorm:"primary_key;column:map_id;type:int;" json:"map_id"`
	//[ 1] description                                    varchar(64)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 64      default: []
	Description string `gorm:"column:description;type:varchar;size:64;" json:"description"`
	//[ 2] scale                                          int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Scale int32 `gorm:"column:scale;type:int;" json:"scale"`
	//[ 3] width                                          int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Width int32 `gorm:"column:width;type:int;" json:"width"`
	//[ 4] height                                         int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Height int32 `gorm:"column:height;type:int;" json:"height"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderMaps) TableName() string {
	return "roomfinder_maps"
}
