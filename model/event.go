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


CREATE TABLE `event` (
  `event` int NOT NULL AUTO_INCREMENT,
  `news` int DEFAULT NULL,
  `kino` int DEFAULT NULL,
  `file` int DEFAULT NULL,
  `title` varchar(100) NOT NULL,
  `description` text NOT NULL,
  `locality` varchar(200) NOT NULL,
  `link` varchar(200) DEFAULT NULL,
  `start` datetime DEFAULT NULL,
  `end` datetime DEFAULT NULL,
  `ticket_group` int DEFAULT '1',
  PRIMARY KEY (`event`),
  KEY `file` (`file`),
  KEY `fkEventGroup` (`ticket_group`),
  KEY `fkNews` (`news`),
  KEY `fkKino` (`kino`),
  FULLTEXT KEY `searchTitle` (`title`),
  CONSTRAINT `fkEventFile` FOREIGN KEY (`file`) REFERENCES `files` (`file`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fkEventGroup` FOREIGN KEY (`ticket_group`) REFERENCES `ticket_group` (`ticket_group`),
  CONSTRAINT `fkKino` FOREIGN KEY (`kino`) REFERENCES `kino` (`kino`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fkNews` FOREIGN KEY (`news`) REFERENCES `news` (`news`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8

JSON Sample
-------------------------------------
{    "locality": "pdltxXwFssOPlVntgDQavLdqL",    "link": "sqkgaonJOrMWuDLYUPgWQBBcA",    "end": "2105-09-13T14:56:25.969353076+01:00",    "event": 83,    "file": 68,    "description": "ALqBGeHWldmEPQLxuHOQXgIaA",    "start": "2254-05-23T05:05:03.991726682+01:00",    "ticket_group": 95,    "news": 90,    "kino": 41,    "title": "SjyYduigvhPjZQkTfmZTTkYhL"}



*/

// Event struct is a row record of the event table in the tca database
type Event struct {
	//[ 0] event                                          int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Event int32 `gorm:"primary_key;AUTO_INCREMENT;column:event;type:int;" json:"event"`
	//[ 1] news                                           int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	News null.Int `gorm:"column:news;type:int;" json:"news"`
	//[ 2] kino                                           int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Kino null.Int `gorm:"column:kino;type:int;" json:"kino"`
	//[ 3] file                                           int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	File null.Int `gorm:"column:file;type:int;" json:"file"`
	//[ 4] title                                          varchar(100)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	Title string `gorm:"column:title;type:varchar;size:100;" json:"title"`
	//[ 5] description                                    text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Description string `gorm:"column:description;type:text;size:65535;" json:"description"`
	//[ 6] locality                                       varchar(200)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 200     default: []
	Locality string `gorm:"column:locality;type:varchar;size:200;" json:"locality"`
	//[ 7] link                                           varchar(200)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 200     default: []
	Link null.String `gorm:"column:link;type:varchar;size:200;" json:"link"`
	//[ 8] start                                          datetime             null: true   primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	Start null.Time `gorm:"column:start;type:datetime;" json:"start"`
	//[ 9] end                                            datetime             null: true   primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	End null.Time `gorm:"column:end;type:datetime;" json:"end"`
	//[10] ticket_group                                   int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: [1]
	TicketGroup null.Int `gorm:"column:ticket_group;type:int;default:1;" json:"ticket_group"`
}

var eventTableInfo = &TableInfo{
	Name: "event",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "event",
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
			GoFieldName:        "Event",
			GoFieldType:        "int32",
			JSONFieldName:      "event",
			ProtobufFieldName:  "event",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "news",
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
			GoFieldName:        "News",
			GoFieldType:        "null.Int",
			JSONFieldName:      "news",
			ProtobufFieldName:  "news",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "kino",
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
			GoFieldName:        "Kino",
			GoFieldType:        "null.Int",
			JSONFieldName:      "kino",
			ProtobufFieldName:  "kino",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
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
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "title",
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
			GoFieldName:        "Title",
			GoFieldType:        "string",
			JSONFieldName:      "title",
			ProtobufFieldName:  "title",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "description",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(65535)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       65535,
			GoFieldName:        "Description",
			GoFieldType:        "string",
			JSONFieldName:      "description",
			ProtobufFieldName:  "description",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "locality",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(200)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       200,
			GoFieldName:        "Locality",
			GoFieldType:        "string",
			JSONFieldName:      "locality",
			ProtobufFieldName:  "locality",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "link",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(200)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       200,
			GoFieldName:        "Link",
			GoFieldType:        "null.String",
			JSONFieldName:      "link",
			ProtobufFieldName:  "link",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "start",
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
			GoFieldName:        "Start",
			GoFieldType:        "null.Time",
			JSONFieldName:      "start",
			ProtobufFieldName:  "start",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        9,
		},

		&ColumnInfo{
			Index:              9,
			Name:               "end",
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
			GoFieldName:        "End",
			GoFieldType:        "null.Time",
			JSONFieldName:      "end",
			ProtobufFieldName:  "end",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        10,
		},

		&ColumnInfo{
			Index:              10,
			Name:               "ticket_group",
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
			GoFieldName:        "TicketGroup",
			GoFieldType:        "null.Int",
			JSONFieldName:      "ticket_group",
			ProtobufFieldName:  "ticket_group",
			ProtobufType:       "int32",
			ProtobufPos:        11,
		},
	},
}

// TableName sets the insert table name for this struct type
func (e *Event) TableName() string {
	return "event"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (e *Event) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (e *Event) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (e *Event) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (e *Event) TableInfo() *TableInfo {
	return eventTableInfo
}
