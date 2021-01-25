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


CREATE TABLE `ticket_payment` (
  `ticket_payment` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `min_amount` int DEFAULT NULL,
  `max_amount` int DEFAULT NULL,
  `config` text NOT NULL,
  PRIMARY KEY (`ticket_payment`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8

JSON Sample
-------------------------------------
{    "ticket_payment": 72,    "name": "CfQyPKskhXXxvWiBaKsInSSgS",    "min_amount": 35,    "max_amount": 11,    "config": "BFcSqBrbkYsvQeATXWwBjaBkC"}



*/

// TicketPayment struct is a row record of the ticket_payment table in the tca database
type TicketPayment struct {
	//[ 0] ticket_payment                                 int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	TicketPayment int32 `gorm:"primary_key;AUTO_INCREMENT;column:ticket_payment;type:int;" json:"ticket_payment"`
	//[ 1] name                                           varchar(50)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 50      default: []
	Name string `gorm:"column:name;type:varchar;size:50;" json:"name"`
	//[ 2] min_amount                                     int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	MinAmount null.Int `gorm:"column:min_amount;type:int;" json:"min_amount"`
	//[ 3] max_amount                                     int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	MaxAmount null.Int `gorm:"column:max_amount;type:int;" json:"max_amount"`
	//[ 4] config                                         text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Config string `gorm:"column:config;type:text;size:65535;" json:"config"`
}

var ticket_paymentTableInfo = &TableInfo{
	Name: "ticket_payment",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "ticket_payment",
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
			GoFieldName:        "TicketPayment",
			GoFieldType:        "int32",
			JSONFieldName:      "ticket_payment",
			ProtobufFieldName:  "ticket_payment",
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
			Name:               "min_amount",
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
			GoFieldName:        "MinAmount",
			GoFieldType:        "null.Int",
			JSONFieldName:      "min_amount",
			ProtobufFieldName:  "min_amount",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "max_amount",
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
			GoFieldName:        "MaxAmount",
			GoFieldType:        "null.Int",
			JSONFieldName:      "max_amount",
			ProtobufFieldName:  "max_amount",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "config",
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
			GoFieldName:        "Config",
			GoFieldType:        "string",
			JSONFieldName:      "config",
			ProtobufFieldName:  "config",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (t *TicketPayment) TableName() string {
	return "ticket_payment"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (t *TicketPayment) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (t *TicketPayment) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (t *TicketPayment) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (t *TicketPayment) TableInfo() *TableInfo {
	return ticket_paymentTableInfo
}
