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


CREATE TABLE `card_option` (
  `card_option` int NOT NULL AUTO_INCREMENT,
  `card` int NOT NULL,
  `text` varchar(2000) DEFAULT '',
  `is_correct_answer` tinyint(1) DEFAULT '0',
  `sort_order` int NOT NULL DEFAULT '0',
  `image` varchar(2000) DEFAULT NULL,
  PRIMARY KEY (`card_option`),
  KEY `card` (`card`),
  CONSTRAINT `card_option_ibfk_1` FOREIGN KEY (`card`) REFERENCES `card` (`card`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "image": "ScOdTOymfXWbUMGTirgBRBWcm",    "card_option": 60,    "card": 57,    "text": "NqVhHqSrFaQgIbDSGqpOaSIFE",    "is_correct_answer": 97,    "sort_order": 8}



*/

// CardOption struct is a row record of the card_option table in the tca database
type CardOption struct {
	//[ 0] card_option                                    int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	CardOption int32 `gorm:"primary_key;AUTO_INCREMENT;column:card_option;type:int;" json:"card_option"`
	//[ 1] card                                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Card int32 `gorm:"column:card;type:int;" json:"card"`
	//[ 2] text                                           varchar(2000)        null: true   primary: false  isArray: false  auto: false  col: varchar         len: 2000    default: []
	Text null.String `gorm:"column:text;type:varchar;size:2000;" json:"text"`
	//[ 3] is_correct_answer                              tinyint              null: true   primary: false  isArray: false  auto: false  col: tinyint         len: -1      default: [0]
	IsCorrectAnswer null.Int `gorm:"column:is_correct_answer;type:tinyint;default:0;" json:"is_correct_answer"`
	//[ 4] sort_order                                     int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	SortOrder int32 `gorm:"column:sort_order;type:int;default:0;" json:"sort_order"`
	//[ 5] image                                          varchar(2000)        null: true   primary: false  isArray: false  auto: false  col: varchar         len: 2000    default: []
	Image null.String `gorm:"column:image;type:varchar;size:2000;" json:"image"`
}

var card_optionTableInfo = &TableInfo{
	Name: "card_option",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "card_option",
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
			GoFieldName:        "CardOption",
			GoFieldType:        "int32",
			JSONFieldName:      "card_option",
			ProtobufFieldName:  "card_option",
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
			IsPrimaryKey:       false,
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

		&ColumnInfo{
			Index:              2,
			Name:               "text",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(2000)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       2000,
			GoFieldName:        "Text",
			GoFieldType:        "null.String",
			JSONFieldName:      "text",
			ProtobufFieldName:  "text",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "is_correct_answer",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "tinyint",
			DatabaseTypePretty: "tinyint",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "tinyint",
			ColumnLength:       -1,
			GoFieldName:        "IsCorrectAnswer",
			GoFieldType:        "null.Int",
			JSONFieldName:      "is_correct_answer",
			ProtobufFieldName:  "is_correct_answer",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "sort_order",
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
			GoFieldName:        "SortOrder",
			GoFieldType:        "int32",
			JSONFieldName:      "sort_order",
			ProtobufFieldName:  "sort_order",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "image",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(2000)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       2000,
			GoFieldName:        "Image",
			GoFieldType:        "null.String",
			JSONFieldName:      "image",
			ProtobufFieldName:  "image",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *CardOption) TableName() string {
	return "card_option"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *CardOption) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *CardOption) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *CardOption) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *CardOption) TableInfo() *TableInfo {
	return card_optionTableInfo
}
