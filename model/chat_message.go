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


CREATE TABLE `chat_message` (
  `message` int NOT NULL AUTO_INCREMENT,
  `member` int NOT NULL,
  `room` int NOT NULL,
  `text` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created` datetime NOT NULL,
  `signature` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`message`),
  KEY `chat_message_b3c09425` (`member`),
  KEY `chat_message_ca20ebca` (`room`),
  CONSTRAINT `chat_message_ibfk_1` FOREIGN KEY (`member`) REFERENCES `member` (`member`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_chat_message_chat_room` FOREIGN KEY (`room`) REFERENCES `chat_room` (`room`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "room": 32,    "text": "PtHbcKsXSVUaYJOODnYoXumKp",    "created": "2269-11-04T00:06:41.451841396+01:00",    "signature": "RkElvpAmbwdByjAoEWRemapNC",    "message": 86,    "member": 96}



*/

// ChatMessage struct is a row record of the chat_message table in the tca database
type ChatMessage struct {
	//[ 0] message                                        int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Message int32 `gorm:"primary_key;AUTO_INCREMENT;column:message;type:int;" json:"message"`
	//[ 1] member                                         int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Member int32 `gorm:"column:member;type:int;" json:"member"`
	//[ 2] room                                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Room int32 `gorm:"column:room;type:int;" json:"room"`
	//[ 3] text                                           text(4294967295)     null: false  primary: false  isArray: false  auto: false  col: text            len: 4294967295 default: []
	Text string `gorm:"column:text;type:text;size:4294967295;" json:"text"`
	//[ 4] created                                        datetime             null: false  primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	Created time.Time `gorm:"column:created;type:datetime;" json:"created"`
	//[ 5] signature                                      text(4294967295)     null: false  primary: false  isArray: false  auto: false  col: text            len: 4294967295 default: []
	Signature string `gorm:"column:signature;type:text;size:4294967295;" json:"signature"`
}

var chat_messageTableInfo = &TableInfo{
	Name: "chat_message",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "message",
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
			GoFieldName:        "Message",
			GoFieldType:        "int32",
			JSONFieldName:      "message",
			ProtobufFieldName:  "message",
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
			Name:               "room",
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
			GoFieldName:        "Room",
			GoFieldType:        "int32",
			JSONFieldName:      "room",
			ProtobufFieldName:  "room",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "text",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(4294967295)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       4294967295,
			GoFieldName:        "Text",
			GoFieldType:        "string",
			JSONFieldName:      "text",
			ProtobufFieldName:  "text",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "created",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "datetime",
			DatabaseTypePretty: "datetime",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "datetime",
			ColumnLength:       -1,
			GoFieldName:        "Created",
			GoFieldType:        "time.Time",
			JSONFieldName:      "created",
			ProtobufFieldName:  "created",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "signature",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(4294967295)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       4294967295,
			GoFieldName:        "Signature",
			GoFieldType:        "string",
			JSONFieldName:      "signature",
			ProtobufFieldName:  "signature",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *ChatMessage) TableName() string {
	return "chat_message"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *ChatMessage) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *ChatMessage) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *ChatMessage) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *ChatMessage) TableInfo() *TableInfo {
	return chat_messageTableInfo
}
