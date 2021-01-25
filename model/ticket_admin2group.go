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


CREATE TABLE `ticket_admin2group` (
  `ticket_admin2group` int NOT NULL AUTO_INCREMENT,
  `ticket_admin` int NOT NULL,
  `ticket_group` int NOT NULL,
  PRIMARY KEY (`ticket_admin2group`),
  KEY `ticket_admin` (`ticket_admin`),
  KEY `ticket_group` (`ticket_group`),
  CONSTRAINT `fkTicketAdmin` FOREIGN KEY (`ticket_admin`) REFERENCES `ticket_admin` (`ticket_admin`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fkTicketGroup` FOREIGN KEY (`ticket_group`) REFERENCES `ticket_group` (`ticket_group`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8

JSON Sample
-------------------------------------
{    "ticket_admin_2_group": 8,    "ticket_admin": 11,    "ticket_group": 37}



*/

// TicketAdmin2group struct is a row record of the ticket_admin2group table in the tca database
type TicketAdmin2group struct {
	//[ 0] ticket_admin2group                             int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	TicketAdmin2group int32 `gorm:"primary_key;AUTO_INCREMENT;column:ticket_admin2group;type:int;" json:"ticket_admin_2_group"`
	//[ 1] ticket_admin                                   int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	TicketAdmin int32 `gorm:"column:ticket_admin;type:int;" json:"ticket_admin"`
	//[ 2] ticket_group                                   int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	TicketGroup int32 `gorm:"column:ticket_group;type:int;" json:"ticket_group"`
}

var ticket_admin2groupTableInfo = &TableInfo{
	Name: "ticket_admin2group",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "ticket_admin2group",
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
			GoFieldName:        "TicketAdmin2group",
			GoFieldType:        "int32",
			JSONFieldName:      "ticket_admin_2_group",
			ProtobufFieldName:  "ticket_admin_2_group",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "ticket_admin",
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
			GoFieldName:        "TicketAdmin",
			GoFieldType:        "int32",
			JSONFieldName:      "ticket_admin",
			ProtobufFieldName:  "ticket_admin",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "ticket_group",
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
			GoFieldName:        "TicketGroup",
			GoFieldType:        "int32",
			JSONFieldName:      "ticket_group",
			ProtobufFieldName:  "ticket_group",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (t *TicketAdmin2group) TableName() string {
	return "ticket_admin2group"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (t *TicketAdmin2group) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (t *TicketAdmin2group) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (t *TicketAdmin2group) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (t *TicketAdmin2group) TableInfo() *TableInfo {
	return ticket_admin2groupTableInfo
}
