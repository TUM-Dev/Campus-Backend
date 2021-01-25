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


CREATE TABLE `card_box` (
  `card_box` int NOT NULL AUTO_INCREMENT,
  `member` int DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `duration` int NOT NULL,
  PRIMARY KEY (`card_box`),
  KEY `member` (`member`),
  CONSTRAINT `card_box_ibfk_1` FOREIGN KEY (`member`) REFERENCES `member` (`member`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "card_box": 5,    "member": 70,    "title": "KweiXJYYlYRqgkjXKcVioppED",    "duration": 73}



*/

// CardBox struct is a row record of the card_box table in the tca database
type CardBox struct {
	//[ 0] card_box                                       int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	CardBox int32 `gorm:"primary_key;AUTO_INCREMENT;column:card_box;type:int;" json:"card_box"`
	//[ 1] member                                         int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Member null.Int `gorm:"column:member;type:int;" json:"member"`
	//[ 2] title                                          varchar(255)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 255     default: []
	Title null.String `gorm:"column:title;type:varchar;size:255;" json:"title"`
	//[ 3] duration                                       int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Duration int32 `gorm:"column:duration;type:int;" json:"duration"`
}

var card_boxTableInfo = &TableInfo{
	Name: "card_box",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "card_box",
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
			GoFieldName:        "CardBox",
			GoFieldType:        "int32",
			JSONFieldName:      "card_box",
			ProtobufFieldName:  "card_box",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "member",
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
			GoFieldName:        "Member",
			GoFieldType:        "null.Int",
			JSONFieldName:      "member",
			ProtobufFieldName:  "member",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
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
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "duration",
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
			GoFieldName:        "Duration",
			GoFieldType:        "int32",
			JSONFieldName:      "duration",
			ProtobufFieldName:  "duration",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *CardBox) TableName() string {
	return "card_box"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *CardBox) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *CardBox) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *CardBox) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *CardBox) TableInfo() *TableInfo {
	return card_boxTableInfo
}
