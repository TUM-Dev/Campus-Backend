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


CREATE TABLE `question2faculty` (
  `question` int NOT NULL,
  `faculty` int NOT NULL,
  PRIMARY KEY (`question`,`faculty`),
  KEY `faculty` (`faculty`),
  CONSTRAINT `question2faculty_ibfk_1` FOREIGN KEY (`question`) REFERENCES `question` (`question`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `question2faculty_ibfk_2` FOREIGN KEY (`faculty`) REFERENCES `faculty` (`faculty`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "question": 19,    "faculty": 22}



*/

// Question2faculty struct is a row record of the question2faculty table in the tca database
type Question2faculty struct {
	//[ 0] question                                       int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	Question int32 `gorm:"primary_key;column:question;type:int;" json:"question"`
	//[ 1] faculty                                        int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	Faculty int32 `gorm:"primary_key;column:faculty;type:int;" json:"faculty"`
}

var question2facultyTableInfo = &TableInfo{
	Name: "question2faculty",
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
			Name:               "faculty",
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
			GoFieldName:        "Faculty",
			GoFieldType:        "int32",
			JSONFieldName:      "faculty",
			ProtobufFieldName:  "faculty",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (q *Question2faculty) TableName() string {
	return "question2faculty"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (q *Question2faculty) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (q *Question2faculty) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (q *Question2faculty) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (q *Question2faculty) TableInfo() *TableInfo {
	return question2facultyTableInfo
}
