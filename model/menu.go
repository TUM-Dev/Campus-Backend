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


CREATE TABLE `menu` (
  `menu` int NOT NULL AUTO_INCREMENT,
  `right` int DEFAULT NULL,
  `parent` int DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `target` enum('_blank','_self','_parent','_top') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '_self',
  `icon` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `class` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `position` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`menu`),
  KEY `right` (`right`),
  KEY `parent` (`parent`),
  CONSTRAINT `menu_ibfk_1` FOREIGN KEY (`right`) REFERENCES `rights` (`right`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `menu_ibfk_2` FOREIGN KEY (`parent`) REFERENCES `menu` (`menu`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "name": "SwgpdoUhDEQfEBVmBQQyxEZXn",    "path": "ZJrqTkwCesSsrFFBISbqYZXFe",    "target": "eFBqxVxYUmkEQWlbtwRTFnwCF",    "icon": "BadHLnXHMnVPODTWbPMWrZoQI",    "menu": 87,    "right": 26,    "parent": 94,    "class": "XxYfbRHBLKlmJNAFRPfisFlme",    "position": 88}



*/

// Menu struct is a row record of the menu table in the tca database
type Menu struct {
	//[ 0] menu                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Menu int32 `gorm:"primary_key;AUTO_INCREMENT;column:menu;type:int;" json:"menu"`
	//[ 1] right                                          int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Right null.Int `gorm:"column:right;type:int;" json:"right"`
	//[ 2] parent                                         int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Parent null.Int `gorm:"column:parent;type:int;" json:"parent"`
	//[ 3] name                                           varchar(255)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 255     default: []
	Name null.String `gorm:"column:name;type:varchar;size:255;" json:"name"`
	//[ 4] path                                           varchar(255)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 255     default: []
	Path null.String `gorm:"column:path;type:varchar;size:255;" json:"path"`
	//[ 5] target                                         char(7)              null: false  primary: false  isArray: false  auto: false  col: char            len: 7       default: [_self]
	Target string `gorm:"column:target;type:char;size:7;default:_self;" json:"target"`
	//[ 6] icon                                           varchar(200)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 200     default: []
	Icon string `gorm:"column:icon;type:varchar;size:200;" json:"icon"`
	//[ 7] class                                          varchar(200)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 200     default: []
	Class string `gorm:"column:class;type:varchar;size:200;" json:"class"`
	//[ 8] position                                       int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	Position int32 `gorm:"column:position;type:int;default:0;" json:"position"`
}

var menuTableInfo = &TableInfo{
	Name: "menu",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "menu",
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
			GoFieldName:        "Menu",
			GoFieldType:        "int32",
			JSONFieldName:      "menu",
			ProtobufFieldName:  "menu",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
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
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "parent",
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
			GoFieldName:        "Parent",
			GoFieldType:        "null.Int",
			JSONFieldName:      "parent",
			ProtobufFieldName:  "parent",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
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
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "path",
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
			GoFieldName:        "Path",
			GoFieldType:        "null.String",
			JSONFieldName:      "path",
			ProtobufFieldName:  "path",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "target",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "char",
			DatabaseTypePretty: "char(7)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "char",
			ColumnLength:       7,
			GoFieldName:        "Target",
			GoFieldType:        "string",
			JSONFieldName:      "target",
			ProtobufFieldName:  "target",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "icon",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(200)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       200,
			GoFieldName:        "Icon",
			GoFieldType:        "string",
			JSONFieldName:      "icon",
			ProtobufFieldName:  "icon",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "class",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(200)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       200,
			GoFieldName:        "Class",
			GoFieldType:        "string",
			JSONFieldName:      "class",
			ProtobufFieldName:  "class",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "position",
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
			GoFieldName:        "Position",
			GoFieldType:        "int32",
			JSONFieldName:      "position",
			ProtobufFieldName:  "position",
			ProtobufType:       "int32",
			ProtobufPos:        9,
		},
	},
}

// TableName sets the insert table name for this struct type
func (m *Menu) TableName() string {
	return "menu"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (m *Menu) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (m *Menu) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (m *Menu) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (m *Menu) TableInfo() *TableInfo {
	return menuTableInfo
}
