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


CREATE TABLE `files` (
  `file` int NOT NULL AUTO_INCREMENT,
  `name` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `path` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `downloads` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`file`)
) ENGINE=InnoDB AUTO_INCREMENT=3707 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "name": "JjpwUBrWtBTEYOglEfqsNTrae",    "path": "KyxaNxtdqtCyauEnxPafahMZd",    "downloads": 93,    "file": 40}



*/

// Files struct is a row record of the files table in the tca database
type Files struct {
	//[ 0] file                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	File int32 `gorm:"primary_key;AUTO_INCREMENT;column:file;type:int;" json:"file"`
	//[ 1] name                                           text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Name string `gorm:"column:name;type:text;size:16777215;" json:"name"`
	//[ 2] path                                           text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Path string `gorm:"column:path;type:text;size:16777215;" json:"path"`
	//[ 3] downloads                                      int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	Downloads int32 `gorm:"column:downloads;type:int;default:0;" json:"downloads"`
}

var filesTableInfo = &TableInfo{
	Name: "files",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "file",
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
			GoFieldName:        "File",
			GoFieldType:        "int32",
			JSONFieldName:      "file",
			ProtobufFieldName:  "file",
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
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "path",
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
			GoFieldName:        "Path",
			GoFieldType:        "string",
			JSONFieldName:      "path",
			ProtobufFieldName:  "path",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "downloads",
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
			GoFieldName:        "Downloads",
			GoFieldType:        "int32",
			JSONFieldName:      "downloads",
			ProtobufFieldName:  "downloads",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (f *Files) TableName() string {
	return "files"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (f *Files) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (f *Files) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (f *Files) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (f *Files) TableInfo() *TableInfo {
	return filesTableInfo
}
