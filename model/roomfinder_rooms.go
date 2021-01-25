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


CREATE TABLE `roomfinder_rooms` (
  `room_id` int NOT NULL,
  `room_code` varchar(32) DEFAULT NULL,
  `building_nr` varchar(8) DEFAULT NULL,
  `arch_id` varchar(16) DEFAULT NULL,
  `info` varchar(64) DEFAULT NULL,
  `address` varchar(128) DEFAULT NULL,
  `purpose_id` int DEFAULT NULL,
  `purpose` varchar(64) DEFAULT NULL,
  `seats` int DEFAULT NULL,
  `utm_zone` varchar(4) DEFAULT NULL,
  `utm_easting` varchar(32) DEFAULT NULL,
  `utm_northing` varchar(32) DEFAULT NULL,
  `unit_id` int DEFAULT NULL,
  `default_map_id` int DEFAULT NULL,
  PRIMARY KEY (`room_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

JSON Sample
-------------------------------------
{    "arch_id": "vuRLwIyKSViLMIKJCFoKqpLKh",    "info": "jxsIXWSpIaeeGxdqyYdtAWHgR",    "seats": 9,    "unit_id": 70,    "utm_zone": "NViVrySIunqVqIFdaCJtUHUXg",    "room_id": 96,    "room_code": "MiblurpABOmbwEvEqJnLATBDT",    "building_nr": "fhPJcZxAHgusfKlajcSDVHhMM",    "address": "aHcARexLMVXOTnRQHeTaLmrPb",    "purpose_id": 42,    "utm_easting": "sJfUWsFSUglBDVQtOuuWnSOIP",    "purpose": "fWqNGEfOucnWjcbWuxsqroWJS",    "utm_northing": "iXZNnQVIrjSxbeJNxiUBoEMfh",    "default_map_id": 47}



*/

// RoomfinderRooms struct is a row record of the roomfinder_rooms table in the tca database
type RoomfinderRooms struct {
	//[ 0] room_id                                        int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	RoomID int32 `gorm:"primary_key;column:room_id;type:int;" json:"room_id"`
	//[ 1] room_code                                      varchar(32)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	RoomCode null.String `gorm:"column:room_code;type:varchar;size:32;" json:"room_code"`
	//[ 2] building_nr                                    varchar(8)           null: true   primary: false  isArray: false  auto: false  col: varchar         len: 8       default: []
	BuildingNr null.String `gorm:"column:building_nr;type:varchar;size:8;" json:"building_nr"`
	//[ 3] arch_id                                        varchar(16)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 16      default: []
	ArchID null.String `gorm:"column:arch_id;type:varchar;size:16;" json:"arch_id"`
	//[ 4] info                                           varchar(64)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 64      default: []
	Info null.String `gorm:"column:info;type:varchar;size:64;" json:"info"`
	//[ 5] address                                        varchar(128)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 128     default: []
	Address null.String `gorm:"column:address;type:varchar;size:128;" json:"address"`
	//[ 6] purpose_id                                     int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	PurposeID null.Int `gorm:"column:purpose_id;type:int;" json:"purpose_id"`
	//[ 7] purpose                                        varchar(64)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 64      default: []
	Purpose null.String `gorm:"column:purpose;type:varchar;size:64;" json:"purpose"`
	//[ 8] seats                                          int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Seats null.Int `gorm:"column:seats;type:int;" json:"seats"`
	//[ 9] utm_zone                                       varchar(4)           null: true   primary: false  isArray: false  auto: false  col: varchar         len: 4       default: []
	UtmZone null.String `gorm:"column:utm_zone;type:varchar;size:4;" json:"utm_zone"`
	//[10] utm_easting                                    varchar(32)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	UtmEasting null.String `gorm:"column:utm_easting;type:varchar;size:32;" json:"utm_easting"`
	//[11] utm_northing                                   varchar(32)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	UtmNorthing null.String `gorm:"column:utm_northing;type:varchar;size:32;" json:"utm_northing"`
	//[12] unit_id                                        int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	UnitID null.Int `gorm:"column:unit_id;type:int;" json:"unit_id"`
	//[13] default_map_id                                 int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	DefaultMapID null.Int `gorm:"column:default_map_id;type:int;" json:"default_map_id"`
}

var roomfinder_roomsTableInfo = &TableInfo{
	Name: "roomfinder_rooms",
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
			Name:               "room_code",
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
			GoFieldName:        "RoomCode",
			GoFieldType:        "null.String",
			JSONFieldName:      "room_code",
			ProtobufFieldName:  "room_code",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "building_nr",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(8)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       8,
			GoFieldName:        "BuildingNr",
			GoFieldType:        "null.String",
			JSONFieldName:      "building_nr",
			ProtobufFieldName:  "building_nr",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "arch_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(16)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       16,
			GoFieldName:        "ArchID",
			GoFieldType:        "null.String",
			JSONFieldName:      "arch_id",
			ProtobufFieldName:  "arch_id",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "info",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(64)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       64,
			GoFieldName:        "Info",
			GoFieldType:        "null.String",
			JSONFieldName:      "info",
			ProtobufFieldName:  "info",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "address",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(128)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       128,
			GoFieldName:        "Address",
			GoFieldType:        "null.String",
			JSONFieldName:      "address",
			ProtobufFieldName:  "address",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "purpose_id",
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
			GoFieldName:        "PurposeID",
			GoFieldType:        "null.Int",
			JSONFieldName:      "purpose_id",
			ProtobufFieldName:  "purpose_id",
			ProtobufType:       "int32",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "purpose",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(64)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       64,
			GoFieldName:        "Purpose",
			GoFieldType:        "null.String",
			JSONFieldName:      "purpose",
			ProtobufFieldName:  "purpose",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "seats",
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
			GoFieldName:        "Seats",
			GoFieldType:        "null.Int",
			JSONFieldName:      "seats",
			ProtobufFieldName:  "seats",
			ProtobufType:       "int32",
			ProtobufPos:        9,
		},

		&ColumnInfo{
			Index:              9,
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
			ProtobufPos:        10,
		},

		&ColumnInfo{
			Index:              10,
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
			ProtobufPos:        11,
		},

		&ColumnInfo{
			Index:              11,
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
			ProtobufPos:        12,
		},

		&ColumnInfo{
			Index:              12,
			Name:               "unit_id",
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
			GoFieldName:        "UnitID",
			GoFieldType:        "null.Int",
			JSONFieldName:      "unit_id",
			ProtobufFieldName:  "unit_id",
			ProtobufType:       "int32",
			ProtobufPos:        13,
		},

		&ColumnInfo{
			Index:              13,
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
			ProtobufPos:        14,
		},
	},
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderRooms) TableName() string {
	return "roomfinder_rooms"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *RoomfinderRooms) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *RoomfinderRooms) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *RoomfinderRooms) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (r *RoomfinderRooms) TableInfo() *TableInfo {
	return roomfinder_roomsTableInfo
}
