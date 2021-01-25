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


CREATE TABLE `ticket_admin` (
  `ticket_admin` int NOT NULL AUTO_INCREMENT,
  `key` text NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `active` tinyint(1) NOT NULL DEFAULT '0',
  `comment` text,
  PRIMARY KEY (`ticket_admin`)
) ENGINE=InnoDB AUTO_INCREMENT=46 DEFAULT CHARSET=utf8

JSON Sample
-------------------------------------
{    "created": "2217-05-02T14:32:22.511363131+01:00",    "active": 51,    "comment": "RxwicduOMAYDVcTvbATpKuFEx",    "ticket_admin": 41,    "key": "pIDKIhXOSFwSjBpjbRskDeXAI"}



*/

// TicketAdmin struct is a row record of the ticket_admin table in the tca database
type TicketAdmin struct {
	//[ 0] ticket_admin                                   int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	TicketAdmin int32 `gorm:"primary_key;AUTO_INCREMENT;column:ticket_admin;type:int;" json:"ticket_admin"`
	//[ 1] key                                            text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Key string `gorm:"column:key;type:text;size:65535;" json:"key"`
	//[ 2] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 3] active                                         tinyint              null: false  primary: false  isArray: false  auto: false  col: tinyint         len: -1      default: [0]
	Active int32 `gorm:"column:active;type:tinyint;default:0;" json:"active"`
	//[ 4] comment                                        text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Comment null.String `gorm:"column:comment;type:text;size:65535;" json:"comment"`
}

var ticket_adminTableInfo = &TableInfo{
	Name: "ticket_admin",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "ticket_admin",
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
			GoFieldName:        "TicketAdmin",
			GoFieldType:        "int32",
			JSONFieldName:      "ticket_admin",
			ProtobufFieldName:  "ticket_admin",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "key",
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
			GoFieldName:        "Key",
			GoFieldType:        "string",
			JSONFieldName:      "key",
			ProtobufFieldName:  "key",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "created",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "timestamp",
			DatabaseTypePretty: "timestamp",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "timestamp",
			ColumnLength:       -1,
			GoFieldName:        "Created",
			GoFieldType:        "time.Time",
			JSONFieldName:      "created",
			ProtobufFieldName:  "created",
			ProtobufType:       "uint64",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "active",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "tinyint",
			DatabaseTypePretty: "tinyint",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "tinyint",
			ColumnLength:       -1,
			GoFieldName:        "Active",
			GoFieldType:        "int32",
			JSONFieldName:      "active",
			ProtobufFieldName:  "active",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "comment",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(65535)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       65535,
			GoFieldName:        "Comment",
			GoFieldType:        "null.String",
			JSONFieldName:      "comment",
			ProtobufFieldName:  "comment",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (t *TicketAdmin) TableName() string {
	return "ticket_admin"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (t *TicketAdmin) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (t *TicketAdmin) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (t *TicketAdmin) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (t *TicketAdmin) TableInfo() *TableInfo {
	return ticket_adminTableInfo
}
