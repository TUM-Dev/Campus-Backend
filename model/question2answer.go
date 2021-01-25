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


CREATE TABLE `question2answer` (
  `question` int NOT NULL,
  `answer` int NOT NULL,
  `member` int NOT NULL,
  UNIQUE KEY `question` (`question`,`member`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "member": 90,    "question": 46,    "answer": 5}


Comments
-------------------------------------
[ 0] Warning table: question2answer does not have a primary key defined, setting col position 1 question as primary key




*/

// Question2answer struct is a row record of the question2answer table in the tca database
type Question2answer struct {
	//[ 0] question                                       int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	Question int32 `gorm:"primary_key;column:question;type:int;" json:"question"`
	//[ 1] answer                                         int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Answer int32 `gorm:"column:answer;type:int;" json:"answer"`
	//[ 2] member                                         int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Member int32 `gorm:"column:member;type:int;" json:"member"`
}

var question2answerTableInfo = &TableInfo{
	Name: "question2answer",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:   0,
			Name:    "question",
			Comment: ``,
			Notes: `Warning table: question2answer does not have a primary key defined, setting col position 1 question as primary key
`,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
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
			Name:               "answer",
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
			GoFieldName:        "Answer",
			GoFieldType:        "int32",
			JSONFieldName:      "answer",
			ProtobufFieldName:  "answer",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
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
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (q *Question2answer) TableName() string {
	return "question2answer"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (q *Question2answer) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (q *Question2answer) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (q *Question2answer) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (q *Question2answer) TableInfo() *TableInfo {
	return question2answerTableInfo
}
