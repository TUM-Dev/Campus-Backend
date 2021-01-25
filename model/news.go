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


CREATE TABLE `news` (
  `news` int NOT NULL AUTO_INCREMENT,
  `date` datetime NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `title` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `src` int NOT NULL,
  `link` varchar(190) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `image` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `file` int DEFAULT NULL,
  PRIMARY KEY (`news`),
  UNIQUE KEY `link` (`link`),
  KEY `src` (`src`),
  KEY `file` (`file`),
  CONSTRAINT `news_ibfk_1` FOREIGN KEY (`src`) REFERENCES `newsSource` (`source`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `news_ibfk_2` FOREIGN KEY (`file`) REFERENCES `files` (`file`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=687621 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "news": 23,    "created": "2286-02-26T08:11:37.184106384+01:00",    "image": "TVJWqXltsBjYdDAkvfSCEUKfN",    "file": 68,    "date": "2156-09-30T09:00:19.355475712+01:00",    "title": "SqBJtFeahNmTSFgJZiZFvetVT",    "description": "YSLSreooZQuGgrPDHOSQvCfIT",    "src": 54,    "link": "gMbOdlANjtwIPKUBSDbPOmqFe"}



*/

// News struct is a row record of the news table in the tca database
type News struct {
	//[ 0] news                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	News int32 `gorm:"primary_key;AUTO_INCREMENT;column:news;type:int;" json:"news"`
	//[ 1] date                                           datetime             null: false  primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	Date time.Time `gorm:"column:date;type:datetime;" json:"date"`
	//[ 2] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 3] title                                          text(255)            null: false  primary: false  isArray: false  auto: false  col: text            len: 255     default: []
	Title string `gorm:"column:title;type:text;size:255;" json:"title"`
	//[ 4] description                                    text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Description string `gorm:"column:description;type:text;size:65535;" json:"description"`
	//[ 5] src                                            int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Src int32 `gorm:"column:src;type:int;" json:"src"`
	//[ 6] link                                           varchar(190)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 190     default: []
	Link string `gorm:"column:link;type:varchar;size:190;" json:"link"`
	//[ 7] image                                          text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Image null.String `gorm:"column:image;type:text;size:65535;" json:"image"`
	//[ 8] file                                           int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	File null.Int `gorm:"column:file;type:int;" json:"file"`
}

var newsTableInfo = &TableInfo{
	Name: "news",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "news",
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
			GoFieldName:        "News",
			GoFieldType:        "int32",
			JSONFieldName:      "news",
			ProtobufFieldName:  "news",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "date",
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
			GoFieldName:        "Date",
			GoFieldType:        "time.Time",
			JSONFieldName:      "date",
			ProtobufFieldName:  "date",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
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
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "title",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       255,
			GoFieldName:        "Title",
			GoFieldType:        "string",
			JSONFieldName:      "title",
			ProtobufFieldName:  "title",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
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
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "src",
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
			GoFieldName:        "Src",
			GoFieldType:        "int32",
			JSONFieldName:      "src",
			ProtobufFieldName:  "src",
			ProtobufType:       "int32",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "link",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(190)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       190,
			GoFieldName:        "Link",
			GoFieldType:        "string",
			JSONFieldName:      "link",
			ProtobufFieldName:  "link",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "image",
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
			GoFieldName:        "Image",
			GoFieldType:        "null.String",
			JSONFieldName:      "image",
			ProtobufFieldName:  "image",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
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
			ProtobufPos:        9,
		},
	},
}

// TableName sets the insert table name for this struct type
func (n *News) TableName() string {
	return "news"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (n *News) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (n *News) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (n *News) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (n *News) TableInfo() *TableInfo {
	return newsTableInfo
}
