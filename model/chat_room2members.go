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


CREATE TABLE `chat_room2members` (
  `room2members` int NOT NULL AUTO_INCREMENT,
  `room` int NOT NULL,
  `member` int NOT NULL,
  PRIMARY KEY (`room2members`),
  UNIQUE KEY `chatroom_id` (`room`,`member`),
  KEY `chat_chatroom_members_29801a33` (`room`),
  KEY `chat_chatroom_members_b3c09425` (`member`),
  CONSTRAINT `chat_room2members_ibfk_2` FOREIGN KEY (`member`) REFERENCES `member` (`member`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_chat_room2members_chat_room` FOREIGN KEY (`room`) REFERENCES `chat_room` (`room`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "room_2_members": 37,    "room": 84,    "member": 64}



*/

// ChatRoom2members struct is a row record of the chat_room2members table in the tca database
type ChatRoom2members struct {
	//[ 0] room2members                                   int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Room2members int32 `gorm:"primary_key;AUTO_INCREMENT;column:room2members;type:int;" json:"room_2_members"`
	//[ 1] room                                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Room int32 `gorm:"column:room;type:int;" json:"room"`
	//[ 2] member                                         int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Member int32 `gorm:"column:member;type:int;" json:"member"`
}

var chat_room2membersTableInfo = &TableInfo{
	Name: "chat_room2members",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "room2members",
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
			GoFieldName:        "Room2members",
			GoFieldType:        "int32",
			JSONFieldName:      "room_2_members",
			ProtobufFieldName:  "room_2_members",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
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
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
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
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *ChatRoom2members) TableName() string {
	return "chat_room2members"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *ChatRoom2members) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *ChatRoom2members) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *ChatRoom2members) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *ChatRoom2members) TableInfo() *TableInfo {
	return chat_room2membersTableInfo
}
