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


CREATE TABLE `notification_type` (
  `type` int NOT NULL AUTO_INCREMENT,
  `name` text NOT NULL,
  `confirmation` enum('true','false') NOT NULL DEFAULT 'false',
  PRIMARY KEY (`type`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8

JSON Sample
-------------------------------------
{    "type": 48,    "name": "yJFyHGUOPwsDCJMKQJuPILkTO",    "confirmation": "fctnfCgFaOXGAaQTYcqPJqyvN"}



*/

// NotificationType struct is a row record of the notification_type table in the tca database
type NotificationType struct {
	//[ 0] type                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Type int32 `gorm:"primary_key;AUTO_INCREMENT;column:type;type:int;" json:"type"`
	//[ 1] name                                           text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Name string `gorm:"column:name;type:text;size:65535;" json:"name"`
	//[ 2] confirmation                                   char(5)              null: false  primary: false  isArray: false  auto: false  col: char            len: 5       default: [false]
	Confirmation string `gorm:"column:confirmation;type:char;size:5;default:false;" json:"confirmation"`
}

var notification_typeTableInfo = &TableInfo{
	Name: "notification_type",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "type",
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
			GoFieldName:        "Type",
			GoFieldType:        "int32",
			JSONFieldName:      "type",
			ProtobufFieldName:  "type",
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
			DatabaseTypePretty: "text(65535)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       65535,
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "confirmation",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "char",
			DatabaseTypePretty: "char(5)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "char",
			ColumnLength:       5,
			GoFieldName:        "Confirmation",
			GoFieldType:        "string",
			JSONFieldName:      "confirmation",
			ProtobufFieldName:  "confirmation",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (n *NotificationType) TableName() string {
	return "notification_type"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (n *NotificationType) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (n *NotificationType) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (n *NotificationType) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (n *NotificationType) TableInfo() *TableInfo {
	return notification_typeTableInfo
}
