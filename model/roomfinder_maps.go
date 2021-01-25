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


CREATE TABLE `roomfinder_maps` (
  `map_id` int NOT NULL,
  `description` varchar(64) NOT NULL,
  `scale` int NOT NULL,
  `width` int NOT NULL,
  `height` int NOT NULL,
  PRIMARY KEY (`map_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

JSON Sample
-------------------------------------
{    "map_id": 47,    "description": "KtwaKRUMVNxLiEOaQtGsLCbwr",    "scale": 44,    "width": 90,    "height": 45}



*/

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

var roomfinder_mapsTableInfo = &TableInfo{
	Name: "roomfinder_maps",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
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
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "description",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(64)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       64,
			GoFieldName:        "Description",
			GoFieldType:        "string",
			JSONFieldName:      "description",
			ProtobufFieldName:  "description",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "scale",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "Scale",
			GoFieldType:        "int32",
			JSONFieldName:      "scale",
			ProtobufFieldName:  "scale",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "width",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "Width",
			GoFieldType:        "int32",
			JSONFieldName:      "width",
			ProtobufFieldName:  "width",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "height",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "Height",
			GoFieldType:        "int32",
			JSONFieldName:      "height",
			ProtobufFieldName:  "height",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderMaps) TableName() string {
	return "roomfinder_maps"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *RoomfinderMaps) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *RoomfinderMaps) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *RoomfinderMaps) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (r *RoomfinderMaps) TableInfo() *TableInfo {
	return roomfinder_mapsTableInfo
}
