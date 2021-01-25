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


CREATE TABLE `update_note` (
  `version_code` int NOT NULL,
  `version_name` text,
  `message` text,
  PRIMARY KEY (`version_code`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "version_code": 46,    "version_name": "uCQixKoaNeueNggJhNonXUXLN",    "message": "IBrvhxEGJMFUAYQGHFZkNQGGc"}



*/

// UpdateNote struct is a row record of the update_note table in the tca database
type UpdateNote struct {
	//[ 0] version_code                                   int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	VersionCode int32 `gorm:"primary_key;column:version_code;type:int;" json:"version_code"`
	//[ 1] version_name                                   text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	VersionName null.String `gorm:"column:version_name;type:text;size:65535;" json:"version_name"`
	//[ 2] message                                        text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Message null.String `gorm:"column:message;type:text;size:65535;" json:"message"`
}

var update_noteTableInfo = &TableInfo{
	Name: "update_note",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "version_code",
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
			GoFieldName:        "VersionCode",
			GoFieldType:        "int32",
			JSONFieldName:      "version_code",
			ProtobufFieldName:  "version_code",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "version_name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(65535)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       65535,
			GoFieldName:        "VersionName",
			GoFieldType:        "null.String",
			JSONFieldName:      "version_name",
			ProtobufFieldName:  "version_name",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "message",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(65535)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       65535,
			GoFieldName:        "Message",
			GoFieldType:        "null.String",
			JSONFieldName:      "message",
			ProtobufFieldName:  "message",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (u *UpdateNote) TableName() string {
	return "update_note"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (u *UpdateNote) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (u *UpdateNote) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (u *UpdateNote) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (u *UpdateNote) TableInfo() *TableInfo {
	return update_noteTableInfo
}
