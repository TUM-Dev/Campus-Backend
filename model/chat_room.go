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


CREATE TABLE `chat_room` (
  `room` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `semester` varchar(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`room`),
  UNIQUE KEY `Index 2` (`semester`,`name`)
) ENGINE=InnoDB AUTO_INCREMENT=1710061 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "name": "sMXulZjQfhoXjRoWNcYSXNXNq",    "semester": "JUjdsnyAAmoQiYbFmBJxCbhKX",    "room": 63}



*/

// ChatRoom struct is a row record of the chat_room table in the tca database
type ChatRoom struct {
	//[ 0] room                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Room int32 `gorm:"primary_key;AUTO_INCREMENT;column:room;type:int;" json:"room"`
	//[ 1] name                                           varchar(100)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	Name string `gorm:"column:name;type:varchar;size:100;" json:"name"`
	//[ 2] semester                                       varchar(3)           null: true   primary: false  isArray: false  auto: false  col: varchar         len: 3       default: []
	Semester null.String `gorm:"column:semester;type:varchar;size:3;" json:"semester"`
}

var chat_roomTableInfo = &TableInfo{
	Name: "chat_room",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "room",
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
			GoFieldName:        "Room",
			GoFieldType:        "int32",
			JSONFieldName:      "room",
			ProtobufFieldName:  "room",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(100)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       100,
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "semester",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(3)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       3,
			GoFieldName:        "Semester",
			GoFieldType:        "null.String",
			JSONFieldName:      "semester",
			ProtobufFieldName:  "semester",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (c *ChatRoom) TableName() string {
	return "chat_room"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (c *ChatRoom) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (c *ChatRoom) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (c *ChatRoom) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (c *ChatRoom) TableInfo() *TableInfo {
	return chat_roomTableInfo
}
