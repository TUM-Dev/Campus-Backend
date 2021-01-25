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


CREATE TABLE `mensaprices` (
  `price` int NOT NULL AUTO_INCREMENT,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `person` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `typeLong` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `typeNumber` int NOT NULL,
  `value` decimal(10,0) NOT NULL,
  PRIMARY KEY (`price`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "created": "2299-12-02T20:12:28.801746883+01:00",    "person": "NByNKwafWkdTuselnmcbMSwCN",    "type": "EJEgmtoJqcJXcwicwSBKexTiY",    "type_long": "PHZKTwDmMJFXBFpDYwjBXhvZL",    "type_number": 54,    "value": 0.815287684933281,    "price": 35}



*/

// Mensaprices struct is a row record of the mensaprices table in the tca database
type Mensaprices struct {
	//[ 0] price                                          int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Price int32 `gorm:"primary_key;AUTO_INCREMENT;column:price;type:int;" json:"price"`
	//[ 1] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 2] person                                         text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Person string `gorm:"column:person;type:text;size:16777215;" json:"person"`
	//[ 3] type                                           text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Type string `gorm:"column:type;type:text;size:16777215;" json:"type"`
	//[ 4] typeLong                                       text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	TypeLong string `gorm:"column:typeLong;type:text;size:16777215;" json:"type_long"`
	//[ 5] typeNumber                                     int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	TypeNumber int32 `gorm:"column:typeNumber;type:int;" json:"type_number"`
	//[ 6] value                                          decimal              null: false  primary: false  isArray: false  auto: false  col: decimal         len: -1      default: []
	Value float64 `gorm:"column:value;type:decimal;" json:"value"`
}

var mensapricesTableInfo = &TableInfo{
	Name: "mensaprices",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "price",
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
			GoFieldName:        "Price",
			GoFieldType:        "int32",
			JSONFieldName:      "price",
			ProtobufFieldName:  "price",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
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
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "person",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(16777215)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       16777215,
			GoFieldName:        "Person",
			GoFieldType:        "string",
			JSONFieldName:      "person",
			ProtobufFieldName:  "person",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "type",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(16777215)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       16777215,
			GoFieldName:        "Type",
			GoFieldType:        "string",
			JSONFieldName:      "type",
			ProtobufFieldName:  "type",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "typeLong",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(16777215)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       16777215,
			GoFieldName:        "TypeLong",
			GoFieldType:        "string",
			JSONFieldName:      "type_long",
			ProtobufFieldName:  "type_long",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "typeNumber",
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
			GoFieldName:        "TypeNumber",
			GoFieldType:        "int32",
			JSONFieldName:      "type_number",
			ProtobufFieldName:  "type_number",
			ProtobufType:       "int32",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "value",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "decimal",
			DatabaseTypePretty: "decimal",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "decimal",
			ColumnLength:       -1,
			GoFieldName:        "Value",
			GoFieldType:        "float64",
			JSONFieldName:      "value",
			ProtobufFieldName:  "value",
			ProtobufType:       "float",
			ProtobufPos:        7,
		},
	},
}

// TableName sets the insert table name for this struct type
func (m *Mensaprices) TableName() string {
	return "mensaprices"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (m *Mensaprices) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (m *Mensaprices) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (m *Mensaprices) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (m *Mensaprices) TableInfo() *TableInfo {
	return mensapricesTableInfo
}
