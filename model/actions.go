package model

import (
	"database/sql"
	"github.com/google/uuid"
	"time"

	"github.com/guregu/null"
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


CREATE TABLE `actions` (
  `action` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `color` varchar(6) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`action`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "description": "QyXsarYCVvUIJYEjFqNZbSsNd",    "color": "okpNjLkoSpUWXkCISVHtiRuXi",    "action": 66,    "name": "bDCulmtplwDOJGqLhrRXUCDtN"}



*/

// Actions struct is a row record of the actions table in the tca database
type Actions struct {
	//[ 0] action                                         int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Action int32 `gorm:"primary_key;AUTO_INCREMENT;column:action;type:int;" json:"action"`
	//[ 1] name                                           varchar(50)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 50      default: []
	Name string `gorm:"column:name;type:varchar;size:50;" json:"name"`
	//[ 2] description                                    text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Description string `gorm:"column:description;type:text;size:16777215;" json:"description"`
	//[ 3] color                                          varchar(6)           null: false  primary: false  isArray: false  auto: false  col: varchar         len: 6       default: []
	Color string `gorm:"column:color;type:varchar;size:6;" json:"color"`
}

var actionsTableInfo = &TableInfo{
	Name: "actions",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "action",
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
			GoFieldName:        "Action",
			GoFieldType:        "int32",
			JSONFieldName:      "action",
			ProtobufFieldName:  "action",
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
			DatabaseTypePretty: "varchar(50)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       50,
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "description",
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
			GoFieldName:        "Description",
			GoFieldType:        "string",
			JSONFieldName:      "description",
			ProtobufFieldName:  "description",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "color",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(6)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       6,
			GoFieldName:        "Color",
			GoFieldType:        "string",
			JSONFieldName:      "color",
			ProtobufFieldName:  "color",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (a *Actions) TableName() string {
	return "actions"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (a *Actions) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (a *Actions) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (a *Actions) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (a *Actions) TableInfo() *TableInfo {
	return actionsTableInfo
}
