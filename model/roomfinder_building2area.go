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


CREATE TABLE `roomfinder_building2area` (
  `area_id` int NOT NULL,
  `building_nr` varchar(8) NOT NULL,
  `campus` char(1) NOT NULL,
  `name` varchar(32) NOT NULL,
  PRIMARY KEY (`building_nr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

JSON Sample
-------------------------------------
{    "area_id": 93,    "building_nr": "lfOqjRhbJBrOmEcWokPiHXKKl",    "campus": "TqJKIYAkSyjgiAZCVcDCMMBIH",    "name": "qIftZpFFXqydRxrxokSIYygTJ"}



*/

// RoomfinderBuilding2area struct is a row record of the roomfinder_building2area table in the tca database
type RoomfinderBuilding2area struct {
	//[ 0] area_id                                        int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	AreaID int32 `gorm:"column:area_id;type:int;" json:"area_id"`
	//[ 1] building_nr                                    varchar(8)           null: false  primary: true   isArray: false  auto: false  col: varchar         len: 8       default: []
	BuildingNr string `gorm:"primary_key;column:building_nr;type:varchar;size:8;" json:"building_nr"`
	//[ 2] campus                                         char(1)              null: false  primary: false  isArray: false  auto: false  col: char            len: 1       default: []
	Campus string `gorm:"column:campus;type:char;size:1;" json:"campus"`
	//[ 3] name                                           varchar(32)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	Name string `gorm:"column:name;type:varchar;size:32;" json:"name"`
}

var roomfinder_building2areaTableInfo = &TableInfo{
	Name: "roomfinder_building2area",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "area_id",
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
			GoFieldName:        "AreaID",
			GoFieldType:        "int32",
			JSONFieldName:      "area_id",
			ProtobufFieldName:  "area_id",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
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
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "campus",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "char",
			DatabaseTypePretty: "char(1)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "char",
			ColumnLength:       1,
			GoFieldName:        "Campus",
			GoFieldType:        "string",
			JSONFieldName:      "campus",
			ProtobufFieldName:  "campus",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(32)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       32,
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderBuilding2area) TableName() string {
	return "roomfinder_building2area"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *RoomfinderBuilding2area) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *RoomfinderBuilding2area) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *RoomfinderBuilding2area) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (r *RoomfinderBuilding2area) TableInfo() *TableInfo {
	return roomfinder_building2areaTableInfo
}
