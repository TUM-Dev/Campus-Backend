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


CREATE TABLE `dish2mensa` (
  `dish2mensa` int NOT NULL AUTO_INCREMENT,
  `mensa` int NOT NULL,
  `dish` int NOT NULL,
  `date` date NOT NULL,
  `created` datetime NOT NULL,
  `modifierd` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`dish2mensa`),
  KEY `dish` (`dish`),
  KEY `mensa` (`mensa`),
  CONSTRAINT `dish2mensa_ibfk_1` FOREIGN KEY (`mensa`) REFERENCES `mensa` (`mensa`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `dish2mensa_ibfk_2` FOREIGN KEY (`dish`) REFERENCES `dish` (`dish`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "mensa": 74,    "dish": 61,    "date": "2242-10-19T18:50:27.261017541+01:00",    "created": "2243-03-11T07:31:10.182397935+01:00",    "modifierd": "2130-03-04T19:17:53.321129401+01:00",    "dish_2_mensa": 73}



*/

// Dish2mensa struct is a row record of the dish2mensa table in the tca database
type Dish2mensa struct {
	//[ 0] dish2mensa                                     int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Dish2mensa int32 `gorm:"primary_key;AUTO_INCREMENT;column:dish2mensa;type:int;" json:"dish_2_mensa"`
	//[ 1] mensa                                          int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Mensa int32 `gorm:"column:mensa;type:int;" json:"mensa"`
	//[ 2] dish                                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Dish int32 `gorm:"column:dish;type:int;" json:"dish"`
	//[ 3] date                                           date                 null: false  primary: false  isArray: false  auto: false  col: date            len: -1      default: []
	Date time.Time `gorm:"column:date;type:date;" json:"date"`
	//[ 4] created                                        datetime             null: false  primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	Created time.Time `gorm:"column:created;type:datetime;" json:"created"`
	//[ 5] modifierd                                      timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [0000-00-00 00:00:00]
	Modifierd time.Time `gorm:"column:modifierd;type:timestamp;default:0000-00-00 00:00:00;" json:"modifierd"`
}

var dish2mensaTableInfo = &TableInfo{
	Name: "dish2mensa",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "dish2mensa",
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
			GoFieldName:        "Dish2mensa",
			GoFieldType:        "int32",
			JSONFieldName:      "dish_2_mensa",
			ProtobufFieldName:  "dish_2_mensa",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "mensa",
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
			GoFieldName:        "Mensa",
			GoFieldType:        "int32",
			JSONFieldName:      "mensa",
			ProtobufFieldName:  "mensa",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "dish",
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
			GoFieldName:        "Dish",
			GoFieldType:        "int32",
			JSONFieldName:      "dish",
			ProtobufFieldName:  "dish",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "date",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "date",
			DatabaseTypePretty: "date",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "date",
			ColumnLength:       -1,
			GoFieldName:        "Date",
			GoFieldType:        "time.Time",
			JSONFieldName:      "date",
			ProtobufFieldName:  "date",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "created",
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
			GoFieldName:        "Created",
			GoFieldType:        "time.Time",
			JSONFieldName:      "created",
			ProtobufFieldName:  "created",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "modifierd",
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
			GoFieldName:        "Modifierd",
			GoFieldType:        "time.Time",
			JSONFieldName:      "modifierd",
			ProtobufFieldName:  "modifierd",
			ProtobufType:       "uint64",
			ProtobufPos:        6,
		},
	},
}

// TableName sets the insert table name for this struct type
func (d *Dish2mensa) TableName() string {
	return "dish2mensa"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *Dish2mensa) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *Dish2mensa) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *Dish2mensa) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *Dish2mensa) TableInfo() *TableInfo {
	return dish2mensaTableInfo
}
