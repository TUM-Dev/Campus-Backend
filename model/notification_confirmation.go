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


CREATE TABLE `notification_confirmation` (
  `notification` int NOT NULL,
  `device` int NOT NULL,
  `sent` tinyint(1) NOT NULL DEFAULT '0',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `received` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`notification`,`device`),
  KEY `device` (`device`),
  CONSTRAINT `notification_confirmation_ibfk_1` FOREIGN KEY (`notification`) REFERENCES `notification` (`notification`),
  CONSTRAINT `notification_confirmation_ibfk_2` FOREIGN KEY (`device`) REFERENCES `devices` (`device`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

JSON Sample
-------------------------------------
{    "notification": 29,    "device": 35,    "sent": 87,    "created": "2189-04-19T14:13:33.541859093+01:00",    "received": "2182-01-18T08:16:46.087892297+01:00"}



*/

// NotificationConfirmation struct is a row record of the notification_confirmation table in the tca database
type NotificationConfirmation struct {
	//[ 0] notification                                   int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	Notification int32 `gorm:"primary_key;column:notification;type:int;" json:"notification"`
	//[ 1] device                                         int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	Device int32 `gorm:"primary_key;column:device;type:int;" json:"device"`
	//[ 2] sent                                           tinyint              null: false  primary: false  isArray: false  auto: false  col: tinyint         len: -1      default: [0]
	Sent int32 `gorm:"column:sent;type:tinyint;default:0;" json:"sent"`
	//[ 3] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 4] received                                       timestamp            null: true   primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: []
	Received null.Time `gorm:"column:received;type:timestamp;" json:"received"`
}

var notification_confirmationTableInfo = &TableInfo{
	Name: "notification_confirmation",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "notification",
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
			GoFieldName:        "Notification",
			GoFieldType:        "int32",
			JSONFieldName:      "notification",
			ProtobufFieldName:  "notification",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "device",
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
			GoFieldName:        "Device",
			GoFieldType:        "int32",
			JSONFieldName:      "device",
			ProtobufFieldName:  "device",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "sent",
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
			GoFieldName:        "Sent",
			GoFieldType:        "int32",
			JSONFieldName:      "sent",
			ProtobufFieldName:  "sent",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
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
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "received",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "timestamp",
			DatabaseTypePretty: "timestamp",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "timestamp",
			ColumnLength:       -1,
			GoFieldName:        "Received",
			GoFieldType:        "null.Time",
			JSONFieldName:      "received",
			ProtobufFieldName:  "received",
			ProtobufType:       "uint64",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (n *NotificationConfirmation) TableName() string {
	return "notification_confirmation"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (n *NotificationConfirmation) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (n *NotificationConfirmation) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (n *NotificationConfirmation) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (n *NotificationConfirmation) TableInfo() *TableInfo {
	return notification_confirmationTableInfo
}
