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


CREATE TABLE `question` (
  `question` int NOT NULL AUTO_INCREMENT,
  `member` int NOT NULL,
  `text` text NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `end` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`question`),
  KEY `member` (`member`)
) ENGINE=InnoDB AUTO_INCREMENT=282 DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "question": 33,    "member": 76,    "text": "UaIgMxOvyvggCCIvyJyAGmvXD",    "created": "2144-05-10T10:06:51.885751533+01:00",    "end": "2218-08-26T01:47:00.373800839+01:00"}



*/

// Question struct is a row record of the question table in the tca database
type Question struct {
	//[ 0] question                                       int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Question int32 `gorm:"primary_key;AUTO_INCREMENT;column:question;type:int;" json:"question"`
	//[ 1] member                                         int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Member int32 `gorm:"column:member;type:int;" json:"member"`
	//[ 2] text                                           text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Text string `gorm:"column:text;type:text;size:65535;" json:"text"`
	//[ 3] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 4] end                                            timestamp            null: true   primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: []
	End null.Time `gorm:"column:end;type:timestamp;" json:"end"`
}

var questionTableInfo = &TableInfo{
	Name: "question",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "question",
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
			GoFieldName:        "Question",
			GoFieldType:        "int32",
			JSONFieldName:      "question",
			ProtobufFieldName:  "question",
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
			Name:               "text",
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
			GoFieldName:        "Text",
			GoFieldType:        "string",
			JSONFieldName:      "text",
			ProtobufFieldName:  "text",
			ProtobufType:       "string",
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
			Name:               "end",
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
			GoFieldName:        "End",
			GoFieldType:        "null.Time",
			JSONFieldName:      "end",
			ProtobufFieldName:  "end",
			ProtobufType:       "uint64",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (q *Question) TableName() string {
	return "question"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (q *Question) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (q *Question) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (q *Question) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (q *Question) TableInfo() *TableInfo {
	return questionTableInfo
}
