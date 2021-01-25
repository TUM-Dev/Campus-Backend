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


CREATE TABLE `newsSource` (
  `source` int NOT NULL AUTO_INCREMENT,
  `title` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `url` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `icon` int DEFAULT NULL,
  `hook` enum('newspread','impulsivHook') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`source`),
  KEY `icon` (`icon`),
  CONSTRAINT `newsSource_ibfk_1` FOREIGN KEY (`icon`) REFERENCES `files` (`file`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "icon": 48,    "hook": "YyMZgcXwpnqCNUHNDGinpdqMe",    "source": 45,    "title": "rkxgTLaMMnqVfFsFSYrBSlrpL",    "url": "sOxGJSyrMnWeLwRmFhJaaNAgq"}



*/

// NewsSource struct is a row record of the newsSource table in the tca database
type NewsSource struct {
	//[ 0] source                                         int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Source int32 `gorm:"primary_key;AUTO_INCREMENT;column:source;type:int;" json:"source"`
	//[ 1] title                                          text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Title string `gorm:"column:title;type:text;size:16777215;" json:"title"`
	//[ 2] url                                            text(16777215)       null: true   primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	URL null.String `gorm:"column:url;type:text;size:16777215;" json:"url"`
	//[ 3] icon                                           int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Icon null.Int `gorm:"column:icon;type:int;" json:"icon"`
	//[ 4] hook                                           char(12)             null: true   primary: false  isArray: false  auto: false  col: char            len: 12      default: []
	Hook null.String `gorm:"column:hook;type:char;size:12;" json:"hook"`
}

var newsSourceTableInfo = &TableInfo{
	Name: "newsSource",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "source",
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
			GoFieldName:        "Source",
			GoFieldType:        "int32",
			JSONFieldName:      "source",
			ProtobufFieldName:  "source",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "title",
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
			GoFieldName:        "Title",
			GoFieldType:        "string",
			JSONFieldName:      "title",
			ProtobufFieldName:  "title",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "url",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(16777215)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       16777215,
			GoFieldName:        "URL",
			GoFieldType:        "null.String",
			JSONFieldName:      "url",
			ProtobufFieldName:  "url",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "icon",
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
			GoFieldName:        "Icon",
			GoFieldType:        "null.Int",
			JSONFieldName:      "icon",
			ProtobufFieldName:  "icon",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "hook",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "char",
			DatabaseTypePretty: "char(12)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "char",
			ColumnLength:       12,
			GoFieldName:        "Hook",
			GoFieldType:        "null.String",
			JSONFieldName:      "hook",
			ProtobufFieldName:  "hook",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (n *NewsSource) TableName() string {
	return "newsSource"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (n *NewsSource) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (n *NewsSource) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (n *NewsSource) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (n *NewsSource) TableInfo() *TableInfo {
	return newsSourceTableInfo
}
