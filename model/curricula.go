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


CREATE TABLE `curricula` (
  `curriculum` int NOT NULL AUTO_INCREMENT,
  `category` enum('bachelor','master') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'bachelor',
  `name` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `url` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`curriculum`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "curriculum": 21,    "category": "oREKpDVcpMsRvNKhfDUPcvXlD",    "name": "FUIPknCkdEDWIHcmjfOYXXKEn",    "url": "SmPnaQwkXcDoJBklIcKYNaOIr"}



*/

// Curricula struct is a row record of the curricula table in the tca database
type Curricula struct {
	//[ 0] curriculum                                     int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Curriculum int32 `gorm:"primary_key;AUTO_INCREMENT;column:curriculum;type:int;" json:"curriculum"`
	//[ 1] category                                       char(8)              null: false  primary: false  isArray: false  auto: false  col: char            len: 8       default: [bachelor]
	Category string `gorm:"column:category;type:char;size:8;default:bachelor;" json:"category"`
	//[ 2] name                                           text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Name string `gorm:"column:name;type:text;size:16777215;" json:"name"`
	//[ 3] url                                            text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	URL string `gorm:"column:url;type:text;size:16777215;" json:"url"`
}

var curriculaTableInfo = &TableInfo{
	Name: "curricula",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "curriculum",
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
			GoFieldName:        "Curriculum",
			GoFieldType:        "int32",
			JSONFieldName:      "curriculum",
			ProtobufFieldName:  "curriculum",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "category",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "char",
			DatabaseTypePretty: "char(8)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "char",
			ColumnLength:       8,
			GoFieldName:        "Category",
			GoFieldType:        "string",
			JSONFieldName:      "category",
			ProtobufFieldName:  "category",
			ProtobufType:       "string",
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
			Name:               "url",
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
			GoFieldName:        "URL",
			GoFieldType:        "string",
			JSONFieldName:      "url",
			ProtobufFieldName:  "url",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *Curricula) TableName() string {
	return "curricula"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *Curricula) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *Curricula) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *Curricula) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *Curricula) TableInfo() *TableInfo {
	return curriculaTableInfo
}
