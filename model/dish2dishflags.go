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


CREATE TABLE `dish2dishflags` (
  `dish2dishflags` int NOT NULL AUTO_INCREMENT,
  `dish` int NOT NULL,
  `flag` int NOT NULL,
  PRIMARY KEY (`dish2dishflags`),
  UNIQUE KEY `dish` (`dish`,`flag`),
  KEY `flag` (`flag`),
  CONSTRAINT `dish2dishflags_ibfk_1` FOREIGN KEY (`dish`) REFERENCES `dish` (`dish`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `dish2dishflags_ibfk_2` FOREIGN KEY (`flag`) REFERENCES `dishflags` (`flag`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "dish_2_dishflags": 54,    "dish": 66,    "flag": 86}



*/

// Dish2dishflags struct is a row record of the dish2dishflags table in the tca database
type Dish2dishflags struct {
	//[ 0] dish2dishflags                                 int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Dish2dishflags int32 `gorm:"primary_key;AUTO_INCREMENT;column:dish2dishflags;type:int;" json:"dish_2_dishflags"`
	//[ 1] dish                                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Dish int32 `gorm:"column:dish;type:int;" json:"dish"`
	//[ 2] flag                                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Flag int32 `gorm:"column:flag;type:int;" json:"flag"`
}

var dish2dishflagsTableInfo = &TableInfo{
	Name: "dish2dishflags",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "dish2dishflags",
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
			GoFieldName:        "Dish2dishflags",
			GoFieldType:        "int32",
			JSONFieldName:      "dish_2_dishflags",
			ProtobufFieldName:  "dish_2_dishflags",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "dish",
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
			GoFieldName:        "Dish",
			GoFieldType:        "int32",
			JSONFieldName:      "dish",
			ProtobufFieldName:  "dish",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "flag",
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
			GoFieldName:        "Flag",
			GoFieldType:        "int32",
			JSONFieldName:      "flag",
			ProtobufFieldName:  "flag",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (d *Dish2dishflags) TableName() string {
	return "dish2dishflags"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *Dish2dishflags) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *Dish2dishflags) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *Dish2dishflags) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *Dish2dishflags) TableInfo() *TableInfo {
	return dish2dishflagsTableInfo
}
