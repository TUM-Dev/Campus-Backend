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


CREATE TABLE `location` (
  `location` int NOT NULL AUTO_INCREMENT,
  `name` text NOT NULL,
  `lon` float(10,6) NOT NULL,
  `lat` float(10,6) NOT NULL,
  `radius` int NOT NULL DEFAULT '1000' COMMENT 'in meters',
  PRIMARY KEY (`location`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8

JSON Sample
-------------------------------------
{    "location": 7,    "name": "SaigytGHfVnqLVIfxWLAjAsjh",    "lon": 0.85090274,    "lat": 0.22363408,    "radius": 64}



*/

// Location struct is a row record of the location table in the tca database
type Location struct {
	//[ 0] location                                       int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Location int32 `gorm:"primary_key;AUTO_INCREMENT;column:location;type:int;" json:"location"`
	//[ 1] name                                           text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Name string `gorm:"column:name;type:text;size:65535;" json:"name"`
	//[ 2] lon                                            float                null: false  primary: false  isArray: false  auto: false  col: float           len: -1      default: []
	Lon float32 `gorm:"column:lon;type:float;" json:"lon"`
	//[ 3] lat                                            float                null: false  primary: false  isArray: false  auto: false  col: float           len: -1      default: []
	Lat float32 `gorm:"column:lat;type:float;" json:"lat"`
	//[ 4] radius                                         int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [1000]
	Radius int32 `gorm:"column:radius;type:int;default:1000;" json:"radius"` // in meters

}

var locationTableInfo = &TableInfo{
	Name: "location",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "location",
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
			GoFieldName:        "Location",
			GoFieldType:        "int32",
			JSONFieldName:      "location",
			ProtobufFieldName:  "location",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(65535)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       65535,
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "lon",
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
			GoFieldName:        "Lon",
			GoFieldType:        "float32",
			JSONFieldName:      "lon",
			ProtobufFieldName:  "lon",
			ProtobufType:       "float",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "lat",
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
			GoFieldName:        "Lat",
			GoFieldType:        "float32",
			JSONFieldName:      "lat",
			ProtobufFieldName:  "lat",
			ProtobufType:       "float",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "radius",
			Comment:            `in meters`,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "Radius",
			GoFieldType:        "int32",
			JSONFieldName:      "radius",
			ProtobufFieldName:  "radius",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (l *Location) TableName() string {
	return "location"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (l *Location) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (l *Location) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (l *Location) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (l *Location) TableInfo() *TableInfo {
	return locationTableInfo
}
