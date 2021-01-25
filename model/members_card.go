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


CREATE TABLE `members_card` (
  `member` int NOT NULL,
  `card` int NOT NULL,
  `card_box` int DEFAULT NULL,
  `last_answered_active_day` int DEFAULT NULL,
  PRIMARY KEY (`member`,`card`),
  KEY `card` (`card`),
  KEY `card_box` (`card_box`),
  CONSTRAINT `members_card_ibfk_1` FOREIGN KEY (`member`) REFERENCES `member` (`member`) ON DELETE CASCADE,
  CONSTRAINT `members_card_ibfk_2` FOREIGN KEY (`card`) REFERENCES `card` (`card`) ON DELETE CASCADE,
  CONSTRAINT `members_card_ibfk_3` FOREIGN KEY (`card_box`) REFERENCES `card_box` (`card_box`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "member": 14,    "card": 19,    "card_box": 15,    "last_answered_active_day": 46}



*/

// MembersCard struct is a row record of the members_card table in the tca database
type MembersCard struct {
	//[ 0] member                                         int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	Member int32 `gorm:"primary_key;column:member;type:int;" json:"member"`
	//[ 1] card                                           int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	Card int32 `gorm:"primary_key;column:card;type:int;" json:"card"`
	//[ 2] card_box                                       int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	CardBox null.Int `gorm:"column:card_box;type:int;" json:"card_box"`
	//[ 3] last_answered_active_day                       int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	LastAnsweredActiveDay null.Int `gorm:"column:last_answered_active_day;type:int;" json:"last_answered_active_day"`
}

var members_cardTableInfo = &TableInfo{
	Name: "members_card",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "member",
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
			GoFieldName:        "Member",
			GoFieldType:        "int32",
			JSONFieldName:      "member",
			ProtobufFieldName:  "member",
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

		&ColumnInfo{
			Index:              2,
			Name:               "card_box",
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
			GoFieldName:        "CardBox",
			GoFieldType:        "null.Int",
			JSONFieldName:      "card_box",
			ProtobufFieldName:  "card_box",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "last_answered_active_day",
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
			GoFieldName:        "LastAnsweredActiveDay",
			GoFieldType:        "null.Int",
			JSONFieldName:      "last_answered_active_day",
			ProtobufFieldName:  "last_answered_active_day",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (m *MembersCard) TableName() string {
	return "members_card"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (m *MembersCard) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (m *MembersCard) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (m *MembersCard) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (m *MembersCard) TableInfo() *TableInfo {
	return members_cardTableInfo
}
