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

// RoomfinderRooms struct is a row record of the roomfinder_rooms table in the tca database
type RoomfinderRooms struct {
	//[ 0] room_id                                        int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	RoomID int32 `gorm:"primary_key;column:room_id;type:int;" json:"room_id"`
	//[ 1] room_code                                      varchar(32)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	RoomCode null.String `gorm:"column:room_code;type:varchar(32);" json:"room_code"`
	//[ 2] building_nr                                    varchar(8)           null: true   primary: false  isArray: false  auto: false  col: varchar         len: 8       default: []
	BuildingNr null.String `gorm:"column:building_nr;type:varchar(8);" json:"building_nr"`
	//[ 3] arch_id                                        varchar(16)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 16      default: []
	ArchID null.String `gorm:"column:arch_id;type:varchar(16);" json:"arch_id"`
	//[ 4] info                                           varchar(64)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 64      default: []
	Info null.String `gorm:"column:info;type:varchar(64);" json:"info"`
	//[ 5] address                                        varchar(128)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 128     default: []
	Address null.String `gorm:"column:address;type:varchar(128);" json:"address"`
	//[ 6] purpose_id                                     int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	PurposeID null.Int `gorm:"column:purpose_id;type:int;" json:"purpose_id"`
	//[ 7] purpose                                        varchar(64)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 64      default: []
	Purpose null.String `gorm:"column:purpose;type:varchar(64);" json:"purpose"`
	//[ 8] seats                                          int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Seats null.Int `gorm:"column:seats;type:int;" json:"seats"`
	//[ 9] utm_zone                                       varchar(4)           null: true   primary: false  isArray: false  auto: false  col: varchar         len: 4       default: []
	UtmZone null.String `gorm:"column:utm_zone;type:varchar(4);" json:"utm_zone"`
	//[10] utm_easting                                    varchar(32)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	UtmEasting null.String `gorm:"column:utm_easting;type:varchar(32);" json:"utm_easting"`
	//[11] utm_northing                                   varchar(32)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	UtmNorthing null.String `gorm:"column:utm_northing;type:varchar(32);" json:"utm_northing"`
	//[12] unit_id                                        int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	UnitID null.Int `gorm:"column:unit_id;type:int;" json:"unit_id"`
	//[13] default_map_id                                 int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	DefaultMapID null.Int `gorm:"column:default_map_id;type:int;" json:"default_map_id"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderRooms) TableName() string {
	return "roomfinder_rooms"
}
