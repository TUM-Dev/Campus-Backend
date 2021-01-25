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


CREATE TABLE `alarm_ban` (
  `ban` int NOT NULL AUTO_INCREMENT,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ip` binary(16) NOT NULL,
  PRIMARY KEY (`ban`),
  UNIQUE KEY `ip` (`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "ban": 78,    "created": "2147-03-13T10:12:46.408694556+01:00",    "ip": "VA=="}



*/

// AlarmBan struct is a row record of the alarm_ban table in the tca database
type AlarmBan struct {
	//[ 0] ban                                            int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Ban int32 `gorm:"primary_key;AUTO_INCREMENT;column:ban;type:int;" json:"ban"`
	//[ 1] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 2] ip                                             binary               null: false  primary: false  isArray: false  auto: false  col: binary          len: -1      default: []
	IP []byte `gorm:"column:ip;type:binary;" json:"ip"`
}

var alarm_banTableInfo = &TableInfo{
	Name: "alarm_ban",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "ban",
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
			GoFieldName:        "Ban",
			GoFieldType:        "int32",
			JSONFieldName:      "ban",
			ProtobufFieldName:  "ban",
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
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (a *AlarmBan) TableName() string {
	return "alarm_ban"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (a *AlarmBan) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (a *AlarmBan) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (a *AlarmBan) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (a *AlarmBan) TableInfo() *TableInfo {
	return alarm_banTableInfo
}
