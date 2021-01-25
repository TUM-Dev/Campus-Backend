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


CREATE TABLE `news_alert` (
  `news_alert` int NOT NULL AUTO_INCREMENT,
  `file` int DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `link` text,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `from` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `to` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`news_alert`),
  UNIQUE KEY `FK_File` (`file`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

JSON Sample
-------------------------------------
{    "created": "2120-07-03T09:50:13.337572003+01:00",    "from": "2076-09-12T16:25:51.584833116+01:00",    "to": "2060-04-09T14:29:16.342329919+01:00",    "news_alert": 51,    "file": 89,    "name": "rybWAxTnGeaKoAtrKkfQNDhLE",    "link": "PSCkKpiLVurQruXVCFthofMlA"}



*/

// NewsAlert struct is a row record of the news_alert table in the tca database
type NewsAlert struct {
	//[ 0] news_alert                                     int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	NewsAlert int32 `gorm:"primary_key;AUTO_INCREMENT;column:news_alert;type:int;" json:"news_alert"`
	//[ 1] file                                           int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	File null.Int `gorm:"column:file;type:int;" json:"file"`
	//[ 2] name                                           varchar(100)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	Name null.String `gorm:"column:name;type:varchar;size:100;" json:"name"`
	//[ 3] link                                           text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Link null.String `gorm:"column:link;type:text;size:65535;" json:"link"`
	//[ 4] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 5] from                                           datetime             null: false  primary: false  isArray: false  auto: false  col: datetime        len: -1      default: [CURRENT_TIMESTAMP]
	From time.Time `gorm:"column:from;type:datetime;default:CURRENT_TIMESTAMP;" json:"from"`
	//[ 6] to                                             datetime             null: false  primary: false  isArray: false  auto: false  col: datetime        len: -1      default: [CURRENT_TIMESTAMP]
	To time.Time `gorm:"column:to;type:datetime;default:CURRENT_TIMESTAMP;" json:"to"`
}

var news_alertTableInfo = &TableInfo{
	Name: "news_alert",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "news_alert",
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
			GoFieldName:        "NewsAlert",
			GoFieldType:        "int32",
			JSONFieldName:      "news_alert",
			ProtobufFieldName:  "news_alert",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "file",
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
			GoFieldName:        "File",
			GoFieldType:        "null.Int",
			JSONFieldName:      "file",
			ProtobufFieldName:  "file",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(100)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       100,
			GoFieldName:        "Name",
			GoFieldType:        "null.String",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "link",
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
			GoFieldName:        "Link",
			GoFieldType:        "null.String",
			JSONFieldName:      "link",
			ProtobufFieldName:  "link",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "created",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "timestamp",
			DatabaseTypePretty: "timestamp",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "timestamp",
			ColumnLength:       -1,
			GoFieldName:        "Created",
			GoFieldType:        "time.Time",
			JSONFieldName:      "created",
			ProtobufFieldName:  "created",
			ProtobufType:       "uint64",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "from",
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
			GoFieldName:        "From",
			GoFieldType:        "time.Time",
			JSONFieldName:      "from",
			ProtobufFieldName:  "from",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "to",
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
			GoFieldName:        "To",
			GoFieldType:        "time.Time",
			JSONFieldName:      "to",
			ProtobufFieldName:  "to",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        7,
		},
	},
}

// TableName sets the insert table name for this struct type
func (n *NewsAlert) TableName() string {
	return "news_alert"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (n *NewsAlert) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (n *NewsAlert) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (n *NewsAlert) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (n *NewsAlert) TableInfo() *TableInfo {
	return news_alertTableInfo
}
