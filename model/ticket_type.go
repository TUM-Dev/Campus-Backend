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


CREATE TABLE `ticket_type` (
  `ticket_type` int NOT NULL AUTO_INCREMENT,
  `event` int NOT NULL,
  `ticket_payment` int NOT NULL,
  `price` double NOT NULL,
  `contingent` int NOT NULL,
  `description` varchar(100) NOT NULL,
  PRIMARY KEY (`ticket_type`),
  KEY `event` (`event`),
  KEY `ticket_payment` (`ticket_payment`),
  CONSTRAINT `fkEvent` FOREIGN KEY (`event`) REFERENCES `event` (`event`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fkPayment` FOREIGN KEY (`ticket_payment`) REFERENCES `ticket_payment` (`ticket_payment`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=57 DEFAULT CHARSET=utf8

JSON Sample
-------------------------------------
{    "event": 45,    "ticket_payment": 88,    "price": 0.6912642549665597,    "contingent": 65,    "description": "bpOSgplagoypjNXhyugwjiCPU",    "ticket_type": 31}



*/

// TicketType struct is a row record of the ticket_type table in the tca database
type TicketType struct {
	//[ 0] ticket_type                                    int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	TicketType int32 `gorm:"primary_key;AUTO_INCREMENT;column:ticket_type;type:int;" json:"ticket_type"`
	//[ 1] event                                          int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Event int32 `gorm:"column:event;type:int;" json:"event"`
	//[ 2] ticket_payment                                 int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	TicketPayment int32 `gorm:"column:ticket_payment;type:int;" json:"ticket_payment"`
	//[ 3] price                                          double               null: false  primary: false  isArray: false  auto: false  col: double          len: -1      default: []
	Price float64 `gorm:"column:price;type:double;" json:"price"`
	//[ 4] contingent                                     int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Contingent int32 `gorm:"column:contingent;type:int;" json:"contingent"`
	//[ 5] description                                    varchar(100)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	Description string `gorm:"column:description;type:varchar;size:100;" json:"description"`
}

var ticket_typeTableInfo = &TableInfo{
	Name: "ticket_type",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "ticket_type",
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
			GoFieldName:        "TicketType",
			GoFieldType:        "int32",
			JSONFieldName:      "ticket_type",
			ProtobufFieldName:  "ticket_type",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "event",
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
			GoFieldName:        "Event",
			GoFieldType:        "int32",
			JSONFieldName:      "event",
			ProtobufFieldName:  "event",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "ticket_payment",
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
			GoFieldName:        "TicketPayment",
			GoFieldType:        "int32",
			JSONFieldName:      "ticket_payment",
			ProtobufFieldName:  "ticket_payment",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "price",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "double",
			DatabaseTypePretty: "double",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "double",
			ColumnLength:       -1,
			GoFieldName:        "Price",
			GoFieldType:        "float64",
			JSONFieldName:      "price",
			ProtobufFieldName:  "price",
			ProtobufType:       "float",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "contingent",
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
			GoFieldName:        "Contingent",
			GoFieldType:        "int32",
			JSONFieldName:      "contingent",
			ProtobufFieldName:  "contingent",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "description",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(100)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       100,
			GoFieldName:        "Description",
			GoFieldType:        "string",
			JSONFieldName:      "description",
			ProtobufFieldName:  "description",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},
	},
}

// TableName sets the insert table name for this struct type
func (t *TicketType) TableName() string {
	return "ticket_type"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (t *TicketType) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (t *TicketType) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (t *TicketType) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (t *TicketType) TableInfo() *TableInfo {
	return ticket_typeTableInfo
}
