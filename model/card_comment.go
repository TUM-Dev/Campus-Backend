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


CREATE TABLE `card_comment` (
  `card_comment` int NOT NULL AUTO_INCREMENT,
  `member` int DEFAULT NULL,
  `card` int NOT NULL,
  `rating` int DEFAULT '0',
  `created_at` date NOT NULL,
  PRIMARY KEY (`card_comment`),
  KEY `member` (`member`),
  KEY `card` (`card`),
  CONSTRAINT `card_comment_ibfk_1` FOREIGN KEY (`member`) REFERENCES `member` (`member`) ON DELETE SET NULL,
  CONSTRAINT `card_comment_ibfk_2` FOREIGN KEY (`card`) REFERENCES `card` (`card`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "card_comment": 70,    "member": 57,    "card": 5,    "rating": 60,    "created_at": "2024-06-19T01:43:15.823011254+02:00"}



*/

// CardComment struct is a row record of the card_comment table in the tca database
type CardComment struct {
	//[ 0] card_comment                                   int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	CardComment int32 `gorm:"primary_key;AUTO_INCREMENT;column:card_comment;type:int;" json:"card_comment"`
	//[ 1] member                                         int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Member null.Int `gorm:"column:member;type:int;" json:"member"`
	//[ 2] card                                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Card int32 `gorm:"column:card;type:int;" json:"card"`
	//[ 3] rating                                         int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	Rating null.Int `gorm:"column:rating;type:int;default:0;" json:"rating"`
	//[ 4] created_at                                     date                 null: false  primary: false  isArray: false  auto: false  col: date            len: -1      default: []
	CreatedAt time.Time `gorm:"column:created_at;type:date;" json:"created_at"`
}

var card_commentTableInfo = &TableInfo{
	Name: "card_comment",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "card_comment",
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
			GoFieldName:        "CardComment",
			GoFieldType:        "int32",
			JSONFieldName:      "card_comment",
			ProtobufFieldName:  "card_comment",
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
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "rating",
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
			GoFieldName:        "Rating",
			GoFieldType:        "null.Int",
			JSONFieldName:      "rating",
			ProtobufFieldName:  "rating",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "created_at",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "date",
			DatabaseTypePretty: "date",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "date",
			ColumnLength:       -1,
			GoFieldName:        "CreatedAt",
			GoFieldType:        "time.Time",
			JSONFieldName:      "created_at",
			ProtobufFieldName:  "created_at",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *CardComment) TableName() string {
	return "card_comment"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *CardComment) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *CardComment) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *CardComment) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *CardComment) TableInfo() *TableInfo {
	return card_commentTableInfo
}
