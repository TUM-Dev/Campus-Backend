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


CREATE TABLE `lecture` (
  `lecture` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`lecture`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "lecture": 72,    "title": "rrZlGSnLpdhOZKMspWtoFGiNL"}



*/

// Lecture struct is a row record of the lecture table in the tca database
type Lecture struct {
	//[ 0] lecture                                        int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Lecture int32 `gorm:"primary_key;AUTO_INCREMENT;column:lecture;type:int;" json:"lecture"`
	//[ 1] title                                          varchar(255)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 255     default: []
	Title null.String `gorm:"column:title;type:varchar;size:255;" json:"title"`
}

var lectureTableInfo = &TableInfo{
	Name: "lecture",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "lecture",
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
			GoFieldName:        "Lecture",
			GoFieldType:        "int32",
			JSONFieldName:      "lecture",
			ProtobufFieldName:  "lecture",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "title",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       255,
			GoFieldName:        "Title",
			GoFieldType:        "null.String",
			JSONFieldName:      "title",
			ProtobufFieldName:  "title",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (l *Lecture) TableName() string {
	return "lecture"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (l *Lecture) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (l *Lecture) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (l *Lecture) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (l *Lecture) TableInfo() *TableInfo {
	return lectureTableInfo
}
