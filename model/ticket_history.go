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


CREATE TABLE `ticket_history` (
  `ticket_history` int NOT NULL AUTO_INCREMENT,
  `member` int NOT NULL,
  `ticket_payment` int DEFAULT NULL,
  `ticket_type` int NOT NULL,
  `purchase` datetime DEFAULT NULL,
  `redemption` datetime DEFAULT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `code` char(128) NOT NULL,
  PRIMARY KEY (`ticket_history`),
  KEY `member` (`member`),
  KEY `ticket_payment` (`ticket_payment`),
  KEY `ticket_type` (`ticket_type`),
  CONSTRAINT `fkMember` FOREIGN KEY (`member`) REFERENCES `member` (`member`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fkTicketPayment` FOREIGN KEY (`ticket_payment`) REFERENCES `ticket_payment` (`ticket_payment`) ON UPDATE CASCADE,
  CONSTRAINT `fkTicketType` FOREIGN KEY (`ticket_type`) REFERENCES `ticket_type` (`ticket_type`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8

JSON Sample
-------------------------------------
{    "ticket_history": 61,    "member": 2,    "ticket_payment": 11,    "ticket_type": 75,    "purchase": "2044-08-25T20:57:14.353466707+01:00",    "redemption": "2205-02-04T08:44:54.325449944+01:00",    "created": "2259-12-20T09:41:23.584462369+01:00",    "code": "EpFaAWELwNEwHfsfLaLorHLov"}



*/

// TicketHistory struct is a row record of the ticket_history table in the tca database
type TicketHistory struct {
	//[ 0] ticket_history                                 int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	TicketHistory int32 `gorm:"primary_key;AUTO_INCREMENT;column:ticket_history;type:int;" json:"ticket_history"`
	//[ 1] member                                         int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Member int32 `gorm:"column:member;type:int;" json:"member"`
	//[ 2] ticket_payment                                 int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	TicketPayment null.Int `gorm:"column:ticket_payment;type:int;" json:"ticket_payment"`
	//[ 3] ticket_type                                    int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	TicketType int32 `gorm:"column:ticket_type;type:int;" json:"ticket_type"`
	//[ 4] purchase                                       datetime             null: true   primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	Purchase null.Time `gorm:"column:purchase;type:datetime;" json:"purchase"`
	//[ 5] redemption                                     datetime             null: true   primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	Redemption null.Time `gorm:"column:redemption;type:datetime;" json:"redemption"`
	//[ 6] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 7] code                                           char(128)            null: false  primary: false  isArray: false  auto: false  col: char            len: 128     default: []
	Code string `gorm:"column:code;type:char;size:128;" json:"code"`
}

var ticket_historyTableInfo = &TableInfo{
	Name: "ticket_history",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "ticket_history",
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
			GoFieldName:        "TicketHistory",
			GoFieldType:        "int32",
			JSONFieldName:      "ticket_history",
			ProtobufFieldName:  "ticket_history",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "member",
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
			GoFieldName:        "Member",
			GoFieldType:        "int32",
			JSONFieldName:      "member",
			ProtobufFieldName:  "member",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "ticket_payment",
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
			GoFieldName:        "TicketPayment",
			GoFieldType:        "null.Int",
			JSONFieldName:      "ticket_payment",
			ProtobufFieldName:  "ticket_payment",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "ticket_type",
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
			GoFieldName:        "TicketType",
			GoFieldType:        "int32",
			JSONFieldName:      "ticket_type",
			ProtobufFieldName:  "ticket_type",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "purchase",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "datetime",
			DatabaseTypePretty: "datetime",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "datetime",
			ColumnLength:       -1,
			GoFieldName:        "Purchase",
			GoFieldType:        "null.Time",
			JSONFieldName:      "purchase",
			ProtobufFieldName:  "purchase",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "redemption",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "datetime",
			DatabaseTypePretty: "datetime",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "datetime",
			ColumnLength:       -1,
			GoFieldName:        "Redemption",
			GoFieldType:        "null.Time",
			JSONFieldName:      "redemption",
			ProtobufFieldName:  "redemption",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
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
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "code",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "char",
			DatabaseTypePretty: "char(128)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "char",
			ColumnLength:       128,
			GoFieldName:        "Code",
			GoFieldType:        "string",
			JSONFieldName:      "code",
			ProtobufFieldName:  "code",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},
	},
}

// TableName sets the insert table name for this struct type
func (t *TicketHistory) TableName() string {
	return "ticket_history"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (t *TicketHistory) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (t *TicketHistory) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (t *TicketHistory) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (t *TicketHistory) TableInfo() *TableInfo {
	return ticket_historyTableInfo
}
