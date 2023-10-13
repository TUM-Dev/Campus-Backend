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

// RoomfinderRooms2maps struct is a row record of the roomfinder_rooms2maps table in the tca database
type RoomfinderRooms2maps struct {
	RoomID int64 `gorm:"primary_key;column:room_id;type:int;" json:"room_id"`
	MapID  int64 `gorm:"primary_key;column:map_id;type:int;" json:"map_id"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderRooms2maps) TableName() string {
	return "roomfinder_rooms2maps"
}
