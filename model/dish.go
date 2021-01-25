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


CREATE TABLE `dish` (
  `dish` int NOT NULL AUTO_INCREMENT,
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`dish`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "dish": 42,    "name": "YKmGODUjHfafyAAekptqhrQlx",    "type": "OctXCyCZTPWRvbGEIXHkHFoSs"}



*/

// Dish struct is a row record of the dish table in the tca database
type Dish struct {
	//[ 0] dish                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Dish int32 `gorm:"primary_key;AUTO_INCREMENT;column:dish;type:int;" json:"dish"`
	//[ 1] name                                           varchar(150)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 150     default: []
	Name string `gorm:"column:name;type:varchar;size:150;" json:"name"`
	//[ 2] type                                           varchar(20)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 20      default: []
	Type string `gorm:"column:type;type:varchar;size:20;" json:"type"`
}

var dishTableInfo = &TableInfo{
	Name: "dish",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "dish",
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
			GoFieldName:        "Dish",
			GoFieldType:        "int32",
			JSONFieldName:      "dish",
			ProtobufFieldName:  "dish",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(150)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       150,
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "type",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(20)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       20,
			GoFieldName:        "Type",
			GoFieldType:        "string",
			JSONFieldName:      "type",
			ProtobufFieldName:  "type",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (d *Dish) TableName() string {
	return "dish"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *Dish) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *Dish) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *Dish) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *Dish) TableInfo() *TableInfo {
	return dishTableInfo
}
