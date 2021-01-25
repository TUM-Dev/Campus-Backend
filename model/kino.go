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


CREATE TABLE `kino` (
  `kino` int NOT NULL AUTO_INCREMENT,
  `date` datetime NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `title` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `year` varchar(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `runtime` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `genre` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `director` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `actors` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `rating` varchar(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `cover` int DEFAULT NULL,
  `trailer` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `link` varchar(190) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`kino`),
  UNIQUE KEY `link` (`link`),
  KEY `cover` (`cover`),
  CONSTRAINT `kino_ibfk_1` FOREIGN KEY (`cover`) REFERENCES `files` (`file`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=179 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "year": "uCTWqcYpDAguKZVManAlfygNe",    "runtime": "CeSukiIZKiQptjidchHWSgpHw",    "director": "ndDKgkPsUdZZRXYNCqmjfQFTd",    "created": "2107-12-26T19:50:15.270222232+01:00",    "title": "UcoOyuAStRkdrjLCqeMFZhlOj",    "date": "2065-12-16T02:50:26.003107106+01:00",    "genre": "MHigYJvkZBtabZPTRwLNlPFMW",    "rating": "fErgDlEpGxwdSmPGeNFNCIoGu",    "cover": 24,    "trailer": "MiUTyAmBDWZdPkkucLeFxlJSL",    "kino": 91,    "actors": "rGJiLBTeoYnQSigvaHbFZHKKm",    "description": "ISMtgSIqOLaJefMkQYeeTWJxd",    "link": "aIUEeyiPukXvbFaIbsYlkWLUt"}



*/

// Kino struct is a row record of the kino table in the tca database
type Kino struct {
	//[ 0] kino                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Kino int32 `gorm:"primary_key;AUTO_INCREMENT;column:kino;type:int;" json:"kino"`
	//[ 1] date                                           datetime             null: false  primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	Date time.Time `gorm:"column:date;type:datetime;" json:"date"`
	//[ 2] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 3] title                                          text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Title string `gorm:"column:title;type:text;size:65535;" json:"title"`
	//[ 4] year                                           varchar(4)           null: false  primary: false  isArray: false  auto: false  col: varchar         len: 4       default: []
	Year string `gorm:"column:year;type:varchar;size:4;" json:"year"`
	//[ 5] runtime                                        varchar(40)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 40      default: []
	Runtime string `gorm:"column:runtime;type:varchar;size:40;" json:"runtime"`
	//[ 6] genre                                          varchar(100)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	Genre string `gorm:"column:genre;type:varchar;size:100;" json:"genre"`
	//[ 7] director                                       text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Director string `gorm:"column:director;type:text;size:65535;" json:"director"`
	//[ 8] actors                                         text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Actors string `gorm:"column:actors;type:text;size:65535;" json:"actors"`
	//[ 9] rating                                         varchar(4)           null: false  primary: false  isArray: false  auto: false  col: varchar         len: 4       default: []
	Rating string `gorm:"column:rating;type:varchar;size:4;" json:"rating"`
	//[10] description                                    text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Description string `gorm:"column:description;type:text;size:65535;" json:"description"`
	//[11] cover                                          int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Cover null.Int `gorm:"column:cover;type:int;" json:"cover"`
	//[12] trailer                                        text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Trailer null.String `gorm:"column:trailer;type:text;size:65535;" json:"trailer"`
	//[13] link                                           varchar(190)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 190     default: []
	Link string `gorm:"column:link;type:varchar;size:190;" json:"link"`
}

var kinoTableInfo = &TableInfo{
	Name: "kino",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "kino",
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
			GoFieldName:        "Kino",
			GoFieldType:        "int32",
			JSONFieldName:      "kino",
			ProtobufFieldName:  "kino",
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
			DatabaseTypePretty: "text(65535)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       65535,
			GoFieldName:        "Title",
			GoFieldType:        "string",
			JSONFieldName:      "title",
			ProtobufFieldName:  "title",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "year",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(4)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       4,
			GoFieldName:        "Year",
			GoFieldType:        "string",
			JSONFieldName:      "year",
			ProtobufFieldName:  "year",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "runtime",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(40)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       40,
			GoFieldName:        "Runtime",
			GoFieldType:        "string",
			JSONFieldName:      "runtime",
			ProtobufFieldName:  "runtime",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "genre",
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
			GoFieldName:        "Genre",
			GoFieldType:        "string",
			JSONFieldName:      "genre",
			ProtobufFieldName:  "genre",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "director",
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
			GoFieldName:        "Director",
			GoFieldType:        "string",
			JSONFieldName:      "director",
			ProtobufFieldName:  "director",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "actors",
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
			GoFieldName:        "Actors",
			GoFieldType:        "string",
			JSONFieldName:      "actors",
			ProtobufFieldName:  "actors",
			ProtobufType:       "string",
			ProtobufPos:        9,
		},

		&ColumnInfo{
			Index:              9,
			Name:               "rating",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(4)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       4,
			GoFieldName:        "Rating",
			GoFieldType:        "string",
			JSONFieldName:      "rating",
			ProtobufFieldName:  "rating",
			ProtobufType:       "string",
			ProtobufPos:        10,
		},

		&ColumnInfo{
			Index:              10,
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
			ProtobufPos:        11,
		},

		&ColumnInfo{
			Index:              11,
			Name:               "cover",
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
			GoFieldName:        "Cover",
			GoFieldType:        "null.Int",
			JSONFieldName:      "cover",
			ProtobufFieldName:  "cover",
			ProtobufType:       "int32",
			ProtobufPos:        12,
		},

		&ColumnInfo{
			Index:              12,
			Name:               "trailer",
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
			GoFieldName:        "Trailer",
			GoFieldType:        "null.String",
			JSONFieldName:      "trailer",
			ProtobufFieldName:  "trailer",
			ProtobufType:       "string",
			ProtobufPos:        13,
		},

		&ColumnInfo{
			Index:              13,
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
			ProtobufPos:        14,
		},
	},
}

// TableName sets the insert table name for this struct type
func (k *Kino) TableName() string {
	return "kino"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (k *Kino) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (k *Kino) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (k *Kino) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (k *Kino) TableInfo() *TableInfo {
	return kinoTableInfo
}
