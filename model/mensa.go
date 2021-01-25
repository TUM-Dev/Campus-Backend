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


CREATE TABLE `mensa` (
  `mensa` int NOT NULL AUTO_INCREMENT,
  `id` int DEFAULT NULL,
  `name` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `address` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `latitude` float(10,6) NOT NULL DEFAULT '0.000000',
  `longitude` float(10,6) NOT NULL DEFAULT '0.000000',
  PRIMARY KEY (`mensa`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "longitude": 0.8188545,    "mensa": 65,    "id": 37,    "name": "RIReTptDNLyvRxaHJNgdiNsow",    "address": "asBBZvNjLQdiNeDUilKlUwMcC",    "latitude": 0.8992467}



*/

// Mensa struct is a row record of the mensa table in the tca database
type Mensa struct {
	//[ 0] mensa                                          int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Mensa int32 `gorm:"primary_key;AUTO_INCREMENT;column:mensa;type:int;" json:"mensa"`
	//[ 1] id                                             int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	ID null.Int `gorm:"column:id;type:int;" json:"id"`
	//[ 2] name                                           text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Name string `gorm:"column:name;type:text;size:16777215;" json:"name"`
	//[ 3] address                                        text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Address string `gorm:"column:address;type:text;size:16777215;" json:"address"`
	//[ 4] latitude                                       float                null: false  primary: false  isArray: false  auto: false  col: float           len: -1      default: [0.000000]
	Latitude float32 `gorm:"column:latitude;type:float;default:0.000000;" json:"latitude"`
	//[ 5] longitude                                      float                null: false  primary: false  isArray: false  auto: false  col: float           len: -1      default: [0.000000]
	Longitude float32 `gorm:"column:longitude;type:float;default:0.000000;" json:"longitude"`
}

var mensaTableInfo = &TableInfo{
	Name: "mensa",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "mensa",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       true,
			IsAutoIncrement:    true,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "Mensa",
			GoFieldType:        "int32",
			JSONFieldName:      "mensa",
			ProtobufFieldName:  "mensa",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "id",
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
			GoFieldName:        "ID",
			GoFieldType:        "null.Int",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(16777215)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       16777215,
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "address",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(16777215)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       16777215,
			GoFieldName:        "Address",
			GoFieldType:        "string",
			JSONFieldName:      "address",
			ProtobufFieldName:  "address",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "latitude",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "float",
			DatabaseTypePretty: "float",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "float",
			ColumnLength:       -1,
			GoFieldName:        "Latitude",
			GoFieldType:        "float32",
			JSONFieldName:      "latitude",
			ProtobufFieldName:  "latitude",
			ProtobufType:       "float",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "longitude",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "float",
			DatabaseTypePretty: "float",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "float",
			ColumnLength:       -1,
			GoFieldName:        "Longitude",
			GoFieldType:        "float32",
			JSONFieldName:      "longitude",
			ProtobufFieldName:  "longitude",
			ProtobufType:       "float",
			ProtobufPos:        6,
		},
	},
}

// TableName sets the insert table name for this struct type
func (m *Mensa) TableName() string {
	return "mensa"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (m *Mensa) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (m *Mensa) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (m *Mensa) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (m *Mensa) TableInfo() *TableInfo {
	return mensaTableInfo
}
