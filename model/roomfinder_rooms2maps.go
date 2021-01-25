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

/*
DB Table Details
-------------------------------------


CREATE TABLE `roomfinder_rooms2maps` (
  `room_id` int NOT NULL,
  `map_id` int NOT NULL,
  PRIMARY KEY (`room_id`,`map_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

JSON Sample
-------------------------------------
{    "room_id": 79,    "map_id": 34}



*/

// RoomfinderRooms2maps struct is a row record of the roomfinder_rooms2maps table in the tca database
type RoomfinderRooms2maps struct {
	//[ 0] room_id                                        int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	RoomID int32 `gorm:"primary_key;column:room_id;type:int;" json:"room_id"`
	//[ 1] map_id                                         int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	MapID int32 `gorm:"primary_key;column:map_id;type:int;" json:"map_id"`
}

var roomfinder_rooms2mapsTableInfo = &TableInfo{
	Name: "roomfinder_rooms2maps",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "room_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "RoomID",
			GoFieldType:        "int32",
			JSONFieldName:      "room_id",
			ProtobufFieldName:  "room_id",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "map_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "MapID",
			GoFieldType:        "int32",
			JSONFieldName:      "map_id",
			ProtobufFieldName:  "map_id",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderRooms2maps) TableName() string {
	return "roomfinder_rooms2maps"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *RoomfinderRooms2maps) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *RoomfinderRooms2maps) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *RoomfinderRooms2maps) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (r *RoomfinderRooms2maps) TableInfo() *TableInfo {
	return roomfinder_rooms2mapsTableInfo
}
