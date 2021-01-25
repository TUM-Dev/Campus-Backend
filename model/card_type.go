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


CREATE TABLE `card_type` (
  `card_type` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`card_type`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "title": "xSQdHcyUPHpPolxlJyVsHaNvd",    "card_type": 8}



*/

// CardType struct is a row record of the card_type table in the tca database
type CardType struct {
	//[ 0] card_type                                      int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	CardType int32 `gorm:"primary_key;AUTO_INCREMENT;column:card_type;type:int;" json:"card_type"`
	//[ 1] title                                          varchar(255)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 255     default: []
	Title null.String `gorm:"column:title;type:varchar;size:255;" json:"title"`
}

var card_typeTableInfo = &TableInfo{
	Name: "card_type",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "card_type",
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
			GoFieldName:        "CardType",
			GoFieldType:        "int32",
			JSONFieldName:      "card_type",
			ProtobufFieldName:  "card_type",
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
func (c *CardType) TableName() string {
	return "card_type"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *CardType) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *CardType) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *CardType) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *CardType) TableInfo() *TableInfo {
	return card_typeTableInfo
}
