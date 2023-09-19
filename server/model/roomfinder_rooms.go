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

// RoomfinderRooms struct is a row record of the roomfinder_rooms table in the tca database
type RoomfinderRooms struct {
	RoomID       int64       `gorm:"primary_key;column:room_id;type:int;" json:"room_id"`
	RoomCode     null.String `gorm:"column:room_code;type:varchar(32);" json:"room_code"`
	BuildingNr   null.String `gorm:"column:building_nr;type:varchar(8);" json:"building_nr"`
	ArchID       null.String `gorm:"column:arch_id;type:varchar(16);" json:"arch_id"`
	Info         null.String `gorm:"column:info;type:varchar(64);" json:"info"`
	Address      null.String `gorm:"column:address;type:varchar(128);" json:"address"`
	PurposeID    null.Int    `gorm:"column:purpose_id;type:int;" json:"purpose_id"`
	Purpose      null.String `gorm:"column:purpose;type:varchar(64);" json:"purpose"`
	Seats        null.Int    `gorm:"column:seats;type:int;" json:"seats"`
	UtmZone      null.String `gorm:"column:utm_zone;type:varchar(4);" json:"utm_zone"`
	UtmEasting   null.String `gorm:"column:utm_easting;type:varchar(32);" json:"utm_easting"`
	UtmNorthing  null.String `gorm:"column:utm_northing;type:varchar(32);" json:"utm_northing"`
	UnitID       null.Int    `gorm:"column:unit_id;type:int;" json:"unit_id"`
	DefaultMapID null.Int    `gorm:"column:default_map_id;type:int;" json:"default_map_id"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderRooms) TableName() string {
	return "roomfinder_rooms"
}
