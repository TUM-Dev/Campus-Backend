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


CREATE TABLE `card2tag` (
  `tag` int NOT NULL,
  `card` int NOT NULL,
  PRIMARY KEY (`tag`,`card`),
  KEY `card` (`card`),
  CONSTRAINT `card2tag_ibfk_1` FOREIGN KEY (`tag`) REFERENCES `tag` (`tag`) ON DELETE CASCADE,
  CONSTRAINT `card2tag_ibfk_2` FOREIGN KEY (`card`) REFERENCES `card` (`card`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "tag": 0,    "card": 40}



*/

// Card2tag struct is a row record of the card2tag table in the tca database
type Card2tag struct {
	//[ 0] tag                                            int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	Tag int32 `gorm:"primary_key;column:tag;type:int;" json:"tag"`
	//[ 1] card                                           int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	Card int32 `gorm:"primary_key;column:card;type:int;" json:"card"`
}

var card2tagTableInfo = &TableInfo{
	Name: "card2tag",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "tag",
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
			GoFieldName:        "Tag",
			GoFieldType:        "int32",
			JSONFieldName:      "tag",
			ProtobufFieldName:  "tag",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "card",
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
			GoFieldName:        "Card",
			GoFieldType:        "int32",
			JSONFieldName:      "card",
			ProtobufFieldName:  "card",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *Card2tag) TableName() string {
	return "card2tag"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *Card2tag) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *Card2tag) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *Card2tag) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *Card2tag) TableInfo() *TableInfo {
	return card2tagTableInfo
}
