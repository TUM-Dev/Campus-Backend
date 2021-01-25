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


CREATE TABLE `roomfinder_buildings` (
  `building_nr` varchar(8) NOT NULL,
  `utm_zone` varchar(4) DEFAULT NULL,
  `utm_easting` varchar(32) DEFAULT NULL,
  `utm_northing` varchar(32) DEFAULT NULL,
  `default_map_id` int DEFAULT NULL,
  PRIMARY KEY (`building_nr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

JSON Sample
-------------------------------------
{    "default_map_id": 36,    "building_nr": "UEFYPFtZmdDAFIRfBDyPfYuZx",    "utm_zone": "jrkFPNAOkIIMtsUPXNvVPOaOi",    "utm_easting": "MrwbNAeyFEThhdBNbfbYtsKJo",    "utm_northing": "UlaBUCXKiDlbEGGtekdwNGIRs"}



*/

// RoomfinderBuildings struct is a row record of the roomfinder_buildings table in the tca database
type RoomfinderBuildings struct {
	//[ 0] building_nr                                    varchar(8)           null: false  primary: true   isArray: false  auto: false  col: varchar         len: 8       default: []
	BuildingNr string `gorm:"primary_key;column:building_nr;type:varchar;size:8;" json:"building_nr"`
	//[ 1] utm_zone                                       varchar(4)           null: true   primary: false  isArray: false  auto: false  col: varchar         len: 4       default: []
	UtmZone null.String `gorm:"column:utm_zone;type:varchar;size:4;" json:"utm_zone"`
	//[ 2] utm_easting                                    varchar(32)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	UtmEasting null.String `gorm:"column:utm_easting;type:varchar;size:32;" json:"utm_easting"`
	//[ 3] utm_northing                                   varchar(32)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	UtmNorthing null.String `gorm:"column:utm_northing;type:varchar;size:32;" json:"utm_northing"`
	//[ 4] default_map_id                                 int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	DefaultMapID null.Int `gorm:"column:default_map_id;type:int;" json:"default_map_id"`
}

var roomfinder_buildingsTableInfo = &TableInfo{
	Name: "roomfinder_buildings",
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
			Name:               "utm_zone",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(4)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       4,
			GoFieldName:        "UtmZone",
			GoFieldType:        "null.String",
			JSONFieldName:      "utm_zone",
			ProtobufFieldName:  "utm_zone",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "utm_easting",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(32)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       32,
			GoFieldName:        "UtmEasting",
			GoFieldType:        "null.String",
			JSONFieldName:      "utm_easting",
			ProtobufFieldName:  "utm_easting",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "utm_northing",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(32)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       32,
			GoFieldName:        "UtmNorthing",
			GoFieldType:        "null.String",
			JSONFieldName:      "utm_northing",
			ProtobufFieldName:  "utm_northing",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "default_map_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "DefaultMapID",
			GoFieldType:        "null.Int",
			JSONFieldName:      "default_map_id",
			ProtobufFieldName:  "default_map_id",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderBuildings) TableName() string {
	return "roomfinder_buildings"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *RoomfinderBuildings) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *RoomfinderBuildings) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *RoomfinderBuildings) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (r *RoomfinderBuildings) TableInfo() *TableInfo {
	return roomfinder_buildingsTableInfo
}
