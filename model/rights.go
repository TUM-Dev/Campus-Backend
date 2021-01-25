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


CREATE TABLE `rights` (
  `right` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `category` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`right`),
  UNIQUE KEY `Unquie` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "right": 45,    "name": "tstRLRjDrCrhBrPMGairWlTfl",    "description": "nKsdKBCSVxJcFmNqLdpGCKsay",    "category": 16}



*/

// Rights struct is a row record of the rights table in the tca database
type Rights struct {
	//[ 0] right                                          int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Right int32 `gorm:"primary_key;AUTO_INCREMENT;column:right;type:int;" json:"right"`
	//[ 1] name                                           varchar(100)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	Name null.String `gorm:"column:name;type:varchar;size:100;" json:"name"`
	//[ 2] description                                    text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Description string `gorm:"column:description;type:text;size:16777215;" json:"description"`
	//[ 3] category                                       int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	Category int32 `gorm:"column:category;type:int;default:0;" json:"category"`
}

var rightsTableInfo = &TableInfo{
	Name: "rights",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "right",
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
			GoFieldName:        "Right",
			GoFieldType:        "int32",
			JSONFieldName:      "right",
			ProtobufFieldName:  "right",
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
			DatabaseTypePretty: "varchar(100)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       100,
			GoFieldName:        "Name",
			GoFieldType:        "null.String",
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
			Name:               "category",
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
			GoFieldName:        "Category",
			GoFieldType:        "int32",
			JSONFieldName:      "category",
			ProtobufFieldName:  "category",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (r *Rights) TableName() string {
	return "rights"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *Rights) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *Rights) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *Rights) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (r *Rights) TableInfo() *TableInfo {
	return rightsTableInfo
}
