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


CREATE TABLE `roomfinder_buildings2maps` (
  `building_nr` varchar(8) NOT NULL,
  `map_id` int NOT NULL,
  PRIMARY KEY (`building_nr`,`map_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

JSON Sample
-------------------------------------
{    "building_nr": "wMMrSoMkSnCYHGuNBWWyPNyqo",    "map_id": 75}



*/

// RoomfinderBuildings2maps struct is a row record of the roomfinder_buildings2maps table in the tca database
type RoomfinderBuildings2maps struct {
	//[ 0] building_nr                                    varchar(8)           null: false  primary: true   isArray: false  auto: false  col: varchar         len: 8       default: []
	BuildingNr string `gorm:"primary_key;column:building_nr;type:varchar;size:8;" json:"building_nr"`
	//[ 1] map_id                                         int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	MapID int32 `gorm:"primary_key;column:map_id;type:int;" json:"map_id"`
}

var roomfinder_buildings2mapsTableInfo = &TableInfo{
	Name: "roomfinder_buildings2maps",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "building_nr",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(8)",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       8,
			GoFieldName:        "BuildingNr",
			GoFieldType:        "string",
			JSONFieldName:      "building_nr",
			ProtobufFieldName:  "building_nr",
			ProtobufType:       "string",
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
func (r *RoomfinderBuildings2maps) TableName() string {
	return "roomfinder_buildings2maps"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *RoomfinderBuildings2maps) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *RoomfinderBuildings2maps) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *RoomfinderBuildings2maps) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (r *RoomfinderBuildings2maps) TableInfo() *TableInfo {
	return roomfinder_buildings2mapsTableInfo
}
