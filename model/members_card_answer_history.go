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


CREATE TABLE `members_card_answer_history` (
  `members_card_answer_history` int NOT NULL AUTO_INCREMENT,
  `member` int NOT NULL,
  `card` int DEFAULT NULL,
  `card_box` int DEFAULT NULL,
  `answer` varchar(2000) DEFAULT NULL,
  `answer_score` float(10,2) DEFAULT '0.00',
  `created_at` date NOT NULL,
  PRIMARY KEY (`members_card_answer_history`),
  KEY `member` (`member`),
  KEY `card` (`card`),
  KEY `card_box` (`card_box`),
  CONSTRAINT `members_card_answer_history_ibfk_1` FOREIGN KEY (`member`) REFERENCES `member` (`member`) ON DELETE CASCADE,
  CONSTRAINT `members_card_answer_history_ibfk_2` FOREIGN KEY (`card`) REFERENCES `card` (`card`) ON DELETE SET NULL,
  CONSTRAINT `members_card_answer_history_ibfk_3` FOREIGN KEY (`card_box`) REFERENCES `card_box` (`card_box`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "card": 9,    "card_box": 98,    "answer": "oqvlJWwRCZkIZVCgEZVUlwIUT",    "answer_score": 0.35324776,    "created_at": "2092-09-29T19:27:46.776417472+01:00",    "members_card_answer_history": 82,    "member": 33}



*/

// MembersCardAnswerHistory struct is a row record of the members_card_answer_history table in the tca database
type MembersCardAnswerHistory struct {
	//[ 0] members_card_answer_history                    int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	MembersCardAnswerHistory int32 `gorm:"primary_key;AUTO_INCREMENT;column:members_card_answer_history;type:int;" json:"members_card_answer_history"`
	//[ 1] member                                         int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Member int32 `gorm:"column:member;type:int;" json:"member"`
	//[ 2] card                                           int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Card null.Int `gorm:"column:card;type:int;" json:"card"`
	//[ 3] card_box                                       int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	CardBox null.Int `gorm:"column:card_box;type:int;" json:"card_box"`
	//[ 4] answer                                         varchar(2000)        null: true   primary: false  isArray: false  auto: false  col: varchar         len: 2000    default: []
	Answer null.String `gorm:"column:answer;type:varchar;size:2000;" json:"answer"`
	//[ 5] answer_score                                   float                null: true   primary: false  isArray: false  auto: false  col: float           len: -1      default: [0.00]
	AnswerScore null.Float `gorm:"column:answer_score;type:float;default:0.00;" json:"answer_score"`
	//[ 6] created_at                                     date                 null: false  primary: false  isArray: false  auto: false  col: date            len: -1      default: []
	CreatedAt time.Time `gorm:"column:created_at;type:date;" json:"created_at"`
}

var members_card_answer_historyTableInfo = &TableInfo{
	Name: "members_card_answer_history",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "members_card_answer_history",
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
			GoFieldName:        "MembersCardAnswerHistory",
			GoFieldType:        "int32",
			JSONFieldName:      "members_card_answer_history",
			ProtobufFieldName:  "members_card_answer_history",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
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
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "card",
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
			GoFieldName:        "Card",
			GoFieldType:        "null.Int",
			JSONFieldName:      "card",
			ProtobufFieldName:  "card",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
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
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "answer",
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
			GoFieldName:        "Answer",
			GoFieldType:        "null.String",
			JSONFieldName:      "answer",
			ProtobufFieldName:  "answer",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "answer_score",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "float",
			DatabaseTypePretty: "float",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "float",
			ColumnLength:       -1,
			GoFieldName:        "AnswerScore",
			GoFieldType:        "null.Float",
			JSONFieldName:      "answer_score",
			ProtobufFieldName:  "answer_score",
			ProtobufType:       "float",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
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
			ProtobufPos:        7,
		},
	},
}

// TableName sets the insert table name for this struct type
func (m *MembersCardAnswerHistory) TableName() string {
	return "members_card_answer_history"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (m *MembersCardAnswerHistory) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (m *MembersCardAnswerHistory) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (m *MembersCardAnswerHistory) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (m *MembersCardAnswerHistory) TableInfo() *TableInfo {
	return members_card_answer_historyTableInfo
}
