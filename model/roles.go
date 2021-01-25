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


CREATE TABLE `roles` (
  `role` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`role`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "role": 67,    "name": "KluURaUwdcJBlKsLhdAwfAlEg",    "description": "KJNYVsyrBvKeEfTfsPKDqfbBV"}



*/

// Roles struct is a row record of the roles table in the tca database
type Roles struct {
	//[ 0] role                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Role int32 `gorm:"primary_key;AUTO_INCREMENT;column:role;type:int;" json:"role"`
	//[ 1] name                                           varchar(50)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 50      default: []
	Name string `gorm:"column:name;type:varchar;size:50;" json:"name"`
	//[ 2] description                                    text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Description string `gorm:"column:description;type:text;size:16777215;" json:"description"`
}

var rolesTableInfo = &TableInfo{
	Name: "roles",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "role",
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
			GoFieldName:        "Role",
			GoFieldType:        "int32",
			JSONFieldName:      "role",
			ProtobufFieldName:  "role",
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
	},
}

// TableName sets the insert table name for this struct type
func (r *Roles) TableName() string {
	return "roles"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *Roles) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *Roles) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *Roles) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (r *Roles) TableInfo() *TableInfo {
	return rolesTableInfo
}
