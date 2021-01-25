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


CREATE TABLE `wifi_measurement` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `date` date NOT NULL,
  `SSID` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `BSSID` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `dBm` int DEFAULT NULL,
  `accuracyInMeters` float NOT NULL,
  `latitude` double NOT NULL,
  `longitude` double NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "date": "2254-08-25T05:54:38.239740438+01:00",    "ssid": "HWXUolQIktaLhCaHnQvuEwfKT",    "bssid": "oWuoArkytqBpyygcLtWrSVhXc",    "d_bm": 1,    "accuracy_in_meters": 0.7338906,    "latitude": 0.1275721524283264,    "longitude": 0.5387681273224639,    "id": 5}


Comments
-------------------------------------
[ 0] column is set for unsigned



*/

// WifiMeasurement struct is a row record of the wifi_measurement table in the tca database
type WifiMeasurement struct {
	//[ 0] id                                             uint                 null: false  primary: true   isArray: false  auto: true   col: uint            len: -1      default: []
	ID uint32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:uint;" json:"id"`
	//[ 1] date                                           date                 null: false  primary: false  isArray: false  auto: false  col: date            len: -1      default: []
	Date time.Time `gorm:"column:date;type:date;" json:"date"`
	//[ 2] SSID                                           varchar(32)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	SSID string `gorm:"column:SSID;type:varchar;size:32;" json:"ssid"`
	//[ 3] BSSID                                          varchar(64)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 64      default: []
	BSSID string `gorm:"column:BSSID;type:varchar;size:64;" json:"bssid"`
	//[ 4] dBm                                            int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	DBm null.Int `gorm:"column:dBm;type:int;" json:"d_bm"`
	//[ 5] accuracyInMeters                               float                null: false  primary: false  isArray: false  auto: false  col: float           len: -1      default: []
	AccuracyInMeters float32 `gorm:"column:accuracyInMeters;type:float;" json:"accuracy_in_meters"`
	//[ 6] latitude                                       double               null: false  primary: false  isArray: false  auto: false  col: double          len: -1      default: []
	Latitude float64 `gorm:"column:latitude;type:double;" json:"latitude"`
	//[ 7] longitude                                      double               null: false  primary: false  isArray: false  auto: false  col: double          len: -1      default: []
	Longitude float64 `gorm:"column:longitude;type:double;" json:"longitude"`
}

var wifi_measurementTableInfo = &TableInfo{
	Name: "wifi_measurement",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
			Comment:            ``,
			Notes:              `column is set for unsigned`,
			Nullable:           false,
			DatabaseTypeName:   "uint",
			DatabaseTypePretty: "uint",
			IsPrimaryKey:       true,
			IsAutoIncrement:    true,
			IsArray:            false,
			ColumnType:         "uint",
			ColumnLength:       -1,
			GoFieldName:        "ID",
			GoFieldType:        "uint32",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "uint32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
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
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "SSID",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(32)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       32,
			GoFieldName:        "SSID",
			GoFieldType:        "string",
			JSONFieldName:      "ssid",
			ProtobufFieldName:  "ssid",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "BSSID",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(64)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       64,
			GoFieldName:        "BSSID",
			GoFieldType:        "string",
			JSONFieldName:      "bssid",
			ProtobufFieldName:  "bssid",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "dBm",
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
			GoFieldName:        "DBm",
			GoFieldType:        "null.Int",
			JSONFieldName:      "d_bm",
			ProtobufFieldName:  "d_bm",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "accuracyInMeters",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "float",
			DatabaseTypePretty: "float",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "float",
			ColumnLength:       -1,
			GoFieldName:        "AccuracyInMeters",
			GoFieldType:        "float32",
			JSONFieldName:      "accuracy_in_meters",
			ProtobufFieldName:  "accuracy_in_meters",
			ProtobufType:       "float",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "latitude",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "double",
			DatabaseTypePretty: "double",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "double",
			ColumnLength:       -1,
			GoFieldName:        "Latitude",
			GoFieldType:        "float64",
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
			Nullable:           false,
			DatabaseTypeName:   "double",
			DatabaseTypePretty: "double",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "double",
			ColumnLength:       -1,
			GoFieldName:        "Longitude",
			GoFieldType:        "float64",
			JSONFieldName:      "longitude",
			ProtobufFieldName:  "longitude",
			ProtobufType:       "float",
			ProtobufPos:        8,
		},
	},
}

// TableName sets the insert table name for this struct type
func (w *WifiMeasurement) TableName() string {
	return "wifi_measurement"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (w *WifiMeasurement) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (w *WifiMeasurement) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (w *WifiMeasurement) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (w *WifiMeasurement) TableInfo() *TableInfo {
	return wifi_measurementTableInfo
}
