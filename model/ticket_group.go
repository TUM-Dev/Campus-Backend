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


CREATE TABLE `ticket_group` (
  `ticket_group` int NOT NULL AUTO_INCREMENT,
  `description` text NOT NULL,
  PRIMARY KEY (`ticket_group`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8

JSON Sample
-------------------------------------
{    "ticket_group": 46,    "description": "qUgtqKKkFFhdHNGAybtmOdYnP"}



*/

// TicketGroup struct is a row record of the ticket_group table in the tca database
type TicketGroup struct {
	//[ 0] ticket_group                                   int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	TicketGroup int32 `gorm:"primary_key;AUTO_INCREMENT;column:ticket_group;type:int;" json:"ticket_group"`
	//[ 1] description                                    text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Description string `gorm:"column:description;type:text;size:65535;" json:"description"`
}

var ticket_groupTableInfo = &TableInfo{
	Name: "ticket_group",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "ticket_group",
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
			GoFieldName:        "TicketGroup",
			GoFieldType:        "int32",
			JSONFieldName:      "ticket_group",
			ProtobufFieldName:  "ticket_group",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "description",
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
			GoFieldName:        "Description",
			GoFieldType:        "string",
			JSONFieldName:      "description",
			ProtobufFieldName:  "description",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (t *TicketGroup) TableName() string {
	return "ticket_group"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (t *TicketGroup) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (t *TicketGroup) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (t *TicketGroup) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (t *TicketGroup) TableInfo() *TableInfo {
	return ticket_groupTableInfo
}
