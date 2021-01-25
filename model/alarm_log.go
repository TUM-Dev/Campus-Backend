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


CREATE TABLE `alarm_log` (
  `alarm` int NOT NULL AUTO_INCREMENT,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `send` int NOT NULL,
  `received` int NOT NULL,
  `test` enum('true','false') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'false',
  `ip` binary(16) NOT NULL,
  PRIMARY KEY (`alarm`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "send": 28,    "received": 47,    "test": "gwvSRSwNLmUPljGwVXxnKCDfM",    "ip": "IRApGAlHCVZaBj4IOD8jXQBcP2ApGSYzUCJTMxFKMTsmAA==",    "alarm": 26,    "created": "2292-09-26T00:21:35.816615293+01:00",    "message": "MwSomudVdJsAjTMDEnFdoevab"}



*/

// AlarmLog struct is a row record of the alarm_log table in the tca database
type AlarmLog struct {
	//[ 0] alarm                                          int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Alarm int32 `gorm:"primary_key;AUTO_INCREMENT;column:alarm;type:int;" json:"alarm"`
	//[ 1] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 2] message                                        text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Message string `gorm:"column:message;type:text;size:65535;" json:"message"`
	//[ 3] send                                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Send int32 `gorm:"column:send;type:int;" json:"send"`
	//[ 4] received                                       int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Received int32 `gorm:"column:received;type:int;" json:"received"`
	//[ 5] test                                           char(5)              null: false  primary: false  isArray: false  auto: false  col: char            len: 5       default: [false]
	Test string `gorm:"column:test;type:char;size:5;default:false;" json:"test"`
	//[ 6] ip                                             binary               null: false  primary: false  isArray: false  auto: false  col: binary          len: -1      default: []
	IP []byte `gorm:"column:ip;type:binary;" json:"ip"`
}

var alarm_logTableInfo = &TableInfo{
	Name: "alarm_log",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "alarm",
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
			GoFieldName:        "Alarm",
			GoFieldType:        "int32",
			JSONFieldName:      "alarm",
			ProtobufFieldName:  "alarm",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
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
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "message",
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
			GoFieldName:        "Message",
			GoFieldType:        "string",
			JSONFieldName:      "message",
			ProtobufFieldName:  "message",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "send",
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
			GoFieldName:        "Send",
			GoFieldType:        "int32",
			JSONFieldName:      "send",
			ProtobufFieldName:  "send",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "received",
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
			GoFieldName:        "Received",
			GoFieldType:        "int32",
			JSONFieldName:      "received",
			ProtobufFieldName:  "received",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "test",
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
			GoFieldName:        "Test",
			GoFieldType:        "string",
			JSONFieldName:      "test",
			ProtobufFieldName:  "test",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "ip",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "binary",
			DatabaseTypePretty: "binary",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "binary",
			ColumnLength:       -1,
			GoFieldName:        "IP",
			GoFieldType:        "[]byte",
			JSONFieldName:      "ip",
			ProtobufFieldName:  "ip",
			ProtobufType:       "bytes",
			ProtobufPos:        7,
		},
	},
}

// TableName sets the insert table name for this struct type
func (a *AlarmLog) TableName() string {
	return "alarm_log"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (a *AlarmLog) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (a *AlarmLog) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (a *AlarmLog) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (a *AlarmLog) TableInfo() *TableInfo {
	return alarm_logTableInfo
}
