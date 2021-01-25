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


CREATE TABLE `feedback` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email_id` text CHARACTER SET utf8 COLLATE utf8_general_ci,
  `receiver` text CHARACTER SET utf8 COLLATE utf8_general_ci,
  `reply_to` text CHARACTER SET utf8 COLLATE utf8_general_ci,
  `feedback` text CHARACTER SET utf8 COLLATE utf8_general_ci,
  `image_count` int DEFAULT NULL,
  `latitude` decimal(11,8) DEFAULT NULL,
  `longitude` decimal(11,8) DEFAULT NULL,
  `timestamp` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1

JSON Sample
-------------------------------------
{    "receiver": "QUXhxbskeBfRompLPbFjdHhTL",    "reply_to": "iPWyTSwUFkQkRoYuqvwMUkXuc",    "image_count": 71,    "latitude": 0.5023100383694792,    "longitude": 0.10047839601973574,    "timestamp": "2166-08-21T13:32:23.73406831+01:00",    "email_id": "rIKcExicTTOyvWmbNvvfltHmy",    "feedback": "kJFtlQGoXCTfwhGZJQAlLDjRc",    "id": 9}



*/

// Feedback struct is a row record of the feedback table in the tca database
type Feedback struct {
	//[ 0] id                                             int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	ID int32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id"`
	//[ 1] email_id                                       text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	EmailID null.String `gorm:"column:email_id;type:text;size:65535;" json:"email_id"`
	//[ 2] receiver                                       text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Receiver null.String `gorm:"column:receiver;type:text;size:65535;" json:"receiver"`
	//[ 3] reply_to                                       text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	ReplyTo null.String `gorm:"column:reply_to;type:text;size:65535;" json:"reply_to"`
	//[ 4] feedback                                       text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Feedback null.String `gorm:"column:feedback;type:text;size:65535;" json:"feedback"`
	//[ 5] image_count                                    int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	ImageCount null.Int `gorm:"column:image_count;type:int;" json:"image_count"`
	//[ 6] latitude                                       decimal              null: true   primary: false  isArray: false  auto: false  col: decimal         len: -1      default: []
	Latitude null.Float `gorm:"column:latitude;type:decimal;" json:"latitude"`
	//[ 7] longitude                                      decimal              null: true   primary: false  isArray: false  auto: false  col: decimal         len: -1      default: []
	Longitude null.Float `gorm:"column:longitude;type:decimal;" json:"longitude"`
	//[ 8] timestamp                                      datetime             null: true   primary: false  isArray: false  auto: false  col: datetime        len: -1      default: [CURRENT_TIMESTAMP]
	Timestamp null.Time `gorm:"column:timestamp;type:datetime;default:CURRENT_TIMESTAMP;" json:"timestamp"`
}

var feedbackTableInfo = &TableInfo{
	Name: "feedback",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
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
			GoFieldName:        "ID",
			GoFieldType:        "int32",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "email_id",
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
			GoFieldName:        "EmailID",
			GoFieldType:        "null.String",
			JSONFieldName:      "email_id",
			ProtobufFieldName:  "email_id",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "receiver",
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
			GoFieldName:        "Receiver",
			GoFieldType:        "null.String",
			JSONFieldName:      "receiver",
			ProtobufFieldName:  "receiver",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "reply_to",
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
			GoFieldName:        "ReplyTo",
			GoFieldType:        "null.String",
			JSONFieldName:      "reply_to",
			ProtobufFieldName:  "reply_to",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "feedback",
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
			GoFieldName:        "Feedback",
			GoFieldType:        "null.String",
			JSONFieldName:      "feedback",
			ProtobufFieldName:  "feedback",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "image_count",
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
			GoFieldName:        "ImageCount",
			GoFieldType:        "null.Int",
			JSONFieldName:      "image_count",
			ProtobufFieldName:  "image_count",
			ProtobufType:       "int32",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "latitude",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "decimal",
			DatabaseTypePretty: "decimal",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "decimal",
			ColumnLength:       -1,
			GoFieldName:        "Latitude",
			GoFieldType:        "null.Float",
			JSONFieldName:      "latitude",
			ProtobufFieldName:  "latitude",
			ProtobufType:       "float",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "longitude",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "decimal",
			DatabaseTypePretty: "decimal",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "decimal",
			ColumnLength:       -1,
			GoFieldName:        "Longitude",
			GoFieldType:        "null.Float",
			JSONFieldName:      "longitude",
			ProtobufFieldName:  "longitude",
			ProtobufType:       "float",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "timestamp",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "datetime",
			DatabaseTypePretty: "datetime",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "datetime",
			ColumnLength:       -1,
			GoFieldName:        "Timestamp",
			GoFieldType:        "null.Time",
			JSONFieldName:      "timestamp",
			ProtobufFieldName:  "timestamp",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        9,
		},
	},
}

// TableName sets the insert table name for this struct type
func (f *Feedback) TableName() string {
	return "feedback"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (f *Feedback) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (f *Feedback) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (f *Feedback) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (f *Feedback) TableInfo() *TableInfo {
	return feedbackTableInfo
}
