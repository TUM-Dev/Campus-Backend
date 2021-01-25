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


CREATE TABLE `log` (
  `log` int NOT NULL AUTO_INCREMENT,
  `time` int NOT NULL,
  `user_executed` int DEFAULT NULL,
  `user_affected` int DEFAULT NULL,
  `action` int DEFAULT NULL,
  `comment` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`log`),
  KEY `user` (`user_executed`),
  KEY `action` (`action`),
  KEY `user_affected` (`user_affected`),
  CONSTRAINT `fkLog2Actions` FOREIGN KEY (`action`) REFERENCES `actions` (`action`) ON DELETE SET NULL ON UPDATE SET NULL,
  CONSTRAINT `fkLog2UsersAf` FOREIGN KEY (`user_affected`) REFERENCES `users` (`user`) ON DELETE SET NULL ON UPDATE SET NULL,
  CONSTRAINT `fkLog2UsersEx` FOREIGN KEY (`user_executed`) REFERENCES `users` (`user`) ON DELETE SET NULL ON UPDATE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "log": 44,    "time": 42,    "user_executed": 45,    "user_affected": 67,    "action": 65,    "comment": "ngeDakENUyEHQhmWLVJoPXPNC"}



*/

// Log struct is a row record of the log table in the tca database
type Log struct {
	//[ 0] log                                            int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Log int32 `gorm:"primary_key;AUTO_INCREMENT;column:log;type:int;" json:"log"`
	//[ 1] time                                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Time int32 `gorm:"column:time;type:int;" json:"time"`
	//[ 2] user_executed                                  int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	UserExecuted null.Int `gorm:"column:user_executed;type:int;" json:"user_executed"`
	//[ 3] user_affected                                  int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	UserAffected null.Int `gorm:"column:user_affected;type:int;" json:"user_affected"`
	//[ 4] action                                         int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Action null.Int `gorm:"column:action;type:int;" json:"action"`
	//[ 5] comment                                        varchar(255)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 255     default: []
	Comment string `gorm:"column:comment;type:varchar;size:255;" json:"comment"`
}

var logTableInfo = &TableInfo{
	Name: "log",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "log",
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
			GoFieldName:        "Log",
			GoFieldType:        "int32",
			JSONFieldName:      "log",
			ProtobufFieldName:  "log",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "time",
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
			GoFieldName:        "Time",
			GoFieldType:        "int32",
			JSONFieldName:      "time",
			ProtobufFieldName:  "time",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "user_executed",
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
			GoFieldName:        "UserExecuted",
			GoFieldType:        "null.Int",
			JSONFieldName:      "user_executed",
			ProtobufFieldName:  "user_executed",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "user_affected",
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
			GoFieldName:        "UserAffected",
			GoFieldType:        "null.Int",
			JSONFieldName:      "user_affected",
			ProtobufFieldName:  "user_affected",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "action",
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
			GoFieldName:        "Action",
			GoFieldType:        "null.Int",
			JSONFieldName:      "action",
			ProtobufFieldName:  "action",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "comment",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       255,
			GoFieldName:        "Comment",
			GoFieldType:        "string",
			JSONFieldName:      "comment",
			ProtobufFieldName:  "comment",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},
	},
}

// TableName sets the insert table name for this struct type
func (l *Log) TableName() string {
	return "log"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (l *Log) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (l *Log) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (l *Log) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (l *Log) TableInfo() *TableInfo {
	return logTableInfo
}
