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


CREATE TABLE `openinghours` (
  `id` int NOT NULL AUTO_INCREMENT,
  `category` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `address` varchar(140) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `room` varchar(140) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `transport_station` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `opening_hours` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `infos` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `url` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `language` varchar(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'de',
  `reference_id` int DEFAULT '-1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=113 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "category": "agYsMYomkhQbtUZslpKgGQJHr",    "name": "SrsjDbgWuBWIPwvYUpPgbWSOy",    "room": "KOMrwnFsNxTrmyCNeesJrOYsU",    "transport_station": "ULoZTVqWSxgSgERgHKkAYxEZu",    "opening_hours": "aTVFxgGuniAfKHhltjZfaXuLn",    "infos": "XAAWwMMVInHGKJZRTolWpHqhu",    "language": "GbhnBPVOBtVjvSMyvgfCZPDaA",    "id": 60,    "url": "ZSdiNIZADDYACGLanOIandHQB",    "reference_id": 77,    "address": "sHvjsvmVTVxybbNxxIHxsWreJ"}



*/

// Openinghours struct is a row record of the openinghours table in the tca database
type Openinghours struct {
	//[ 0] id                                             int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	ID int32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id"`
	//[ 1] category                                       varchar(20)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 20      default: []
	Category string `gorm:"column:category;type:varchar;size:20;" json:"category"`
	//[ 2] name                                           varchar(60)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 60      default: []
	Name string `gorm:"column:name;type:varchar;size:60;" json:"name"`
	//[ 3] address                                        varchar(140)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 140     default: []
	Address string `gorm:"column:address;type:varchar;size:140;" json:"address"`
	//[ 4] room                                           varchar(140)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 140     default: []
	Room null.String `gorm:"column:room;type:varchar;size:140;" json:"room"`
	//[ 5] transport_station                              varchar(150)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 150     default: []
	TransportStation null.String `gorm:"column:transport_station;type:varchar;size:150;" json:"transport_station"`
	//[ 6] opening_hours                                  varchar(300)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 300     default: []
	OpeningHours null.String `gorm:"column:opening_hours;type:varchar;size:300;" json:"opening_hours"`
	//[ 7] infos                                          varchar(500)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 500     default: []
	Infos null.String `gorm:"column:infos;type:varchar;size:500;" json:"infos"`
	//[ 8] url                                            varchar(300)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 300     default: []
	URL string `gorm:"column:url;type:varchar;size:300;" json:"url"`
	//[ 9] language                                       varchar(2)           null: true   primary: false  isArray: false  auto: false  col: varchar         len: 2       default: [de]
	Language null.String `gorm:"column:language;type:varchar;size:2;default:de;" json:"language"`
	//[10] reference_id                                   int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: [-1]
	ReferenceID null.Int `gorm:"column:reference_id;type:int;default:-1;" json:"reference_id"`
}

var openinghoursTableInfo = &TableInfo{
	Name: "openinghours",
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
			Name:               "category",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(20)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       20,
			GoFieldName:        "Category",
			GoFieldType:        "string",
			JSONFieldName:      "category",
			ProtobufFieldName:  "category",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(60)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       60,
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "address",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(140)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       140,
			GoFieldName:        "Address",
			GoFieldType:        "string",
			JSONFieldName:      "address",
			ProtobufFieldName:  "address",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "room",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(140)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       140,
			GoFieldName:        "Room",
			GoFieldType:        "null.String",
			JSONFieldName:      "room",
			ProtobufFieldName:  "room",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "transport_station",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(150)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       150,
			GoFieldName:        "TransportStation",
			GoFieldType:        "null.String",
			JSONFieldName:      "transport_station",
			ProtobufFieldName:  "transport_station",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "opening_hours",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(300)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       300,
			GoFieldName:        "OpeningHours",
			GoFieldType:        "null.String",
			JSONFieldName:      "opening_hours",
			ProtobufFieldName:  "opening_hours",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "infos",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(500)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       500,
			GoFieldName:        "Infos",
			GoFieldType:        "null.String",
			JSONFieldName:      "infos",
			ProtobufFieldName:  "infos",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "url",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(300)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       300,
			GoFieldName:        "URL",
			GoFieldType:        "string",
			JSONFieldName:      "url",
			ProtobufFieldName:  "url",
			ProtobufType:       "string",
			ProtobufPos:        9,
		},

		&ColumnInfo{
			Index:              9,
			Name:               "language",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(2)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       2,
			GoFieldName:        "Language",
			GoFieldType:        "null.String",
			JSONFieldName:      "language",
			ProtobufFieldName:  "language",
			ProtobufType:       "string",
			ProtobufPos:        10,
		},

		&ColumnInfo{
			Index:              10,
			Name:               "reference_id",
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
			GoFieldName:        "ReferenceID",
			GoFieldType:        "null.Int",
			JSONFieldName:      "reference_id",
			ProtobufFieldName:  "reference_id",
			ProtobufType:       "int32",
			ProtobufPos:        11,
		},
	},
}

// TableName sets the insert table name for this struct type
func (o *Openinghours) TableName() string {
	return "openinghours"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (o *Openinghours) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (o *Openinghours) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (o *Openinghours) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (o *Openinghours) TableInfo() *TableInfo {
	return openinghoursTableInfo
}
