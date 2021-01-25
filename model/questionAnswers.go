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


CREATE TABLE `questionAnswers` (
  `answer` int NOT NULL AUTO_INCREMENT,
  `text` text NOT NULL,
  PRIMARY KEY (`answer`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "answer": 51,    "text": "vSkNLoBxpXLdOYwRIvOPnJqPj"}



*/

// QuestionAnswers struct is a row record of the questionAnswers table in the tca database
type QuestionAnswers struct {
	//[ 0] answer                                         int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Answer int32 `gorm:"primary_key;AUTO_INCREMENT;column:answer;type:int;" json:"answer"`
	//[ 1] text                                           text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Text string `gorm:"column:text;type:text;size:65535;" json:"text"`
}

var questionAnswersTableInfo = &TableInfo{
	Name: "questionAnswers",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "answer",
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
			GoFieldName:        "Answer",
			GoFieldType:        "int32",
			JSONFieldName:      "answer",
			ProtobufFieldName:  "answer",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
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
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (q *QuestionAnswers) TableName() string {
	return "questionAnswers"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (q *QuestionAnswers) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (q *QuestionAnswers) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (q *QuestionAnswers) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (q *QuestionAnswers) TableInfo() *TableInfo {
	return questionAnswersTableInfo
}
