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


CREATE TABLE `modules` (
  `module` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `right` int DEFAULT NULL,
  PRIMARY KEY (`module`),
  KEY `module_right` (`right`),
  CONSTRAINT `fkMod2Rights` FOREIGN KEY (`right`) REFERENCES `rights` (`right`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "module": 49,    "name": "pSOEEaREqdNngQbdEZNmKjGZU",    "right": 86}



*/

// Modules struct is a row record of the modules table in the tca database
type Modules struct {
	//[ 0] module                                         int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Module int32 `gorm:"primary_key;AUTO_INCREMENT;column:module;type:int;" json:"module"`
	//[ 1] name                                           varchar(255)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 255     default: []
	Name null.String `gorm:"column:name;type:varchar;size:255;" json:"name"`
	//[ 2] right                                          int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Right null.Int `gorm:"column:right;type:int;" json:"right"`
}

var modulesTableInfo = &TableInfo{
	Name: "modules",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "module",
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
			GoFieldName:        "Module",
			GoFieldType:        "int32",
			JSONFieldName:      "module",
			ProtobufFieldName:  "module",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       255,
			GoFieldName:        "Name",
			GoFieldType:        "null.String",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "right",
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
			GoFieldName:        "Right",
			GoFieldType:        "null.Int",
			JSONFieldName:      "right",
			ProtobufFieldName:  "right",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (m *Modules) TableName() string {
	return "modules"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (m *Modules) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (m *Modules) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (m *Modules) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (m *Modules) TableInfo() *TableInfo {
	return modulesTableInfo
}
