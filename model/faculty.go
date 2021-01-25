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


CREATE TABLE `faculty` (
  `faculty` int NOT NULL AUTO_INCREMENT,
  `name` varchar(150) NOT NULL,
  PRIMARY KEY (`faculty`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

JSON Sample
-------------------------------------
{    "faculty": 52,    "name": "yXHJFtSfbmOPFkcibJpELmXsm"}



*/

// Faculty struct is a row record of the faculty table in the tca database
type Faculty struct {
	//[ 0] faculty                                        int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Faculty int32 `gorm:"primary_key;AUTO_INCREMENT;column:faculty;type:int;" json:"faculty"`
	//[ 1] name                                           varchar(150)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 150     default: []
	Name string `gorm:"column:name;type:varchar;size:150;" json:"name"`
}

var facultyTableInfo = &TableInfo{
	Name: "faculty",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "faculty",
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
			GoFieldName:        "Faculty",
			GoFieldType:        "int32",
			JSONFieldName:      "faculty",
			ProtobufFieldName:  "faculty",
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
	},
}

// TableName sets the insert table name for this struct type
func (f *Faculty) TableName() string {
	return "faculty"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (f *Faculty) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (f *Faculty) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (f *Faculty) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (f *Faculty) TableInfo() *TableInfo {
	return facultyTableInfo
}
