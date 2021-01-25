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


CREATE TABLE `sessions` (
  `session` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `access` int unsigned DEFAULT NULL,
  `data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`session`),
  UNIQUE KEY `session` (`session`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "session": "pFQuVZUkDpWBpOnZqkgMolfAQ",    "access": 25,    "data": "GpNpUQlOCTCqOoHimtmnlymtn"}


Comments
-------------------------------------
[ 1] column is set for unsigned



*/

// Sessions struct is a row record of the sessions table in the tca database
type Sessions struct {
	//[ 0] session                                        varchar(255)         null: false  primary: true   isArray: false  auto: false  col: varchar         len: 255     default: []
	Session string `gorm:"primary_key;column:session;type:varchar;size:255;" json:"session"`
	//[ 1] access                                         uint                 null: true   primary: false  isArray: false  auto: false  col: uint            len: -1      default: []
	Access null.Int `gorm:"column:access;type:uint;" json:"access"`
	//[ 2] data                                           text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Data null.String `gorm:"column:data;type:text;size:65535;" json:"data"`
}

var sessionsTableInfo = &TableInfo{
	Name: "sessions",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "session",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(255)",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       255,
			GoFieldName:        "Session",
			GoFieldType:        "string",
			JSONFieldName:      "session",
			ProtobufFieldName:  "session",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "access",
			Comment:            ``,
			Notes:              `column is set for unsigned`,
			Nullable:           true,
			DatabaseTypeName:   "uint",
			DatabaseTypePretty: "uint",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "uint",
			ColumnLength:       -1,
			GoFieldName:        "Access",
			GoFieldType:        "null.Int",
			JSONFieldName:      "access",
			ProtobufFieldName:  "access",
			ProtobufType:       "uint32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "data",
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
			GoFieldName:        "Data",
			GoFieldType:        "null.String",
			JSONFieldName:      "data",
			ProtobufFieldName:  "data",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (s *Sessions) TableName() string {
	return "sessions"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (s *Sessions) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (s *Sessions) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (s *Sessions) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (s *Sessions) TableInfo() *TableInfo {
	return sessionsTableInfo
}
