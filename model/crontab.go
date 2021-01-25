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


CREATE TABLE `crontab` (
  `cron` int NOT NULL AUTO_INCREMENT,
  `interval` int NOT NULL DEFAULT '7200',
  `lastRun` int NOT NULL DEFAULT '0',
  `type` enum('news','mensa','chat','kino','roomfinder','ticketsale','alarm') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `id` int DEFAULT NULL,
  UNIQUE KEY `cron` (`cron`)
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "cron": 75,    "interval": 10,    "last_run": 91,    "type": "HhpqCUTtLxKQBaTXFFnSrtGTs",    "id": 69}


Comments
-------------------------------------
[ 0] Warning table: crontab does not have a primary key defined, setting col position 1 cron as primary key




*/

// Crontab struct is a row record of the crontab table in the tca database
type Crontab struct {
	//[ 0] cron                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Cron int32 `gorm:"primary_key;AUTO_INCREMENT;column:cron;type:int;" json:"cron"`
	//[ 1] interval                                       int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [7200]
	Interval int32 `gorm:"column:interval;type:int;default:7200;" json:"interval"`
	//[ 2] lastRun                                        int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	LastRun int32 `gorm:"column:lastRun;type:int;default:0;" json:"last_run"`
	//[ 3] type                                           char(10)             null: true   primary: false  isArray: false  auto: false  col: char            len: 10      default: []
	Type null.String `gorm:"column:type;type:char;size:10;" json:"type"`
	//[ 4] id                                             int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	ID null.Int `gorm:"column:id;type:int;" json:"id"`
}

var crontabTableInfo = &TableInfo{
	Name: "crontab",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:   0,
			Name:    "cron",
			Comment: ``,
			Notes: `Warning table: crontab does not have a primary key defined, setting col position 1 cron as primary key
`,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       true,
			IsAutoIncrement:    true,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "Cron",
			GoFieldType:        "int32",
			JSONFieldName:      "cron",
			ProtobufFieldName:  "cron",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "interval",
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
			GoFieldName:        "Interval",
			GoFieldType:        "int32",
			JSONFieldName:      "interval",
			ProtobufFieldName:  "interval",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "lastRun",
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
			GoFieldName:        "LastRun",
			GoFieldType:        "int32",
			JSONFieldName:      "last_run",
			ProtobufFieldName:  "last_run",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "type",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "char",
			DatabaseTypePretty: "char(10)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "char",
			ColumnLength:       10,
			GoFieldName:        "Type",
			GoFieldType:        "null.String",
			JSONFieldName:      "type",
			ProtobufFieldName:  "type",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "id",
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
			GoFieldName:        "ID",
			GoFieldType:        "null.Int",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *Crontab) TableName() string {
	return "crontab"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *Crontab) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *Crontab) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *Crontab) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *Crontab) TableInfo() *TableInfo {
	return crontabTableInfo
}
