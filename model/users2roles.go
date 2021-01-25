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


CREATE TABLE `users2roles` (
  `user` int NOT NULL,
  `role` int NOT NULL,
  PRIMARY KEY (`user`,`role`),
  KEY `fkUser2RolesRole` (`role`),
  CONSTRAINT `fkUser2RolesRole` FOREIGN KEY (`role`) REFERENCES `roles` (`role`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fkUser2RolesUser` FOREIGN KEY (`user`) REFERENCES `users` (`user`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "user": 97,    "role": 4}



*/

// Users2roles struct is a row record of the users2roles table in the tca database
type Users2roles struct {
	//[ 0] user                                           int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	User int32 `gorm:"primary_key;column:user;type:int;" json:"user"`
	//[ 1] role                                           int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	Role int32 `gorm:"primary_key;column:role;type:int;" json:"role"`
}

var users2rolesTableInfo = &TableInfo{
	Name: "users2roles",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "user",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "User",
			GoFieldType:        "int32",
			JSONFieldName:      "user",
			ProtobufFieldName:  "user",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "role",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "Role",
			GoFieldType:        "int32",
			JSONFieldName:      "role",
			ProtobufFieldName:  "role",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (u *Users2roles) TableName() string {
	return "users2roles"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (u *Users2roles) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (u *Users2roles) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (u *Users2roles) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (u *Users2roles) TableInfo() *TableInfo {
	return users2rolesTableInfo
}
