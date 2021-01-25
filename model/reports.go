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


CREATE TABLE `reports` (
  `report` int NOT NULL AUTO_INCREMENT,
  `device` int DEFAULT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `fixed` enum('true','false') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'false',
  `issue` int DEFAULT NULL,
  `stacktrace` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `stacktraceGroup` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `log` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `package` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `packageVersion` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `packageVersionCode` int NOT NULL DEFAULT '-1',
  `model` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `osVersion` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `networkWifi` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `networkMobile` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `gps` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `screenWidth` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `screenHeight` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `screenOrientation` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `screenDpi` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`report`),
  KEY `device` (`device`),
  CONSTRAINT `reports_ibfk_3` FOREIGN KEY (`device`) REFERENCES `devices` (`device`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "device": 36,    "issue": 51,    "screen_height": "XDmZDdHaVJwARJNfisCWRaQeQ",    "screen_dpi": "HWXFAwsCvhhZQDTIWiteVQUgw",    "report": 76,    "created": "2060-12-12T06:24:58.466861194+01:00",    "stacktrace_group": "IpGnOkKZETxWfXPZIuSHORIFZ",    "log": "KJIxZHTyJyhZCKjHBDuktBXFb",    "package_version": "RYaYKvoaJIxxmfOVyeWUtvlZk",    "network_mobile": "UvjpvovGbPKnCrZNyJWMDVRBQ",    "screen_orientation": "QRIMJLAwQJeVunyyZOytwUbgt",    "package": "LTXbuSkOYAkYgfDhEmwufRNjj",    "package_version_code": 32,    "model": "maUGEoPqmbdWtAatqSoeuIbNX",    "gps": "rDwvRALrEyxBwRjNFtKYANyiv",    "fixed": "KqyCkECaZBZdoucUBwktGBafq",    "stacktrace": "CuuOJcBPskGDdSMPtTFObOJIu",    "os_version": "tytRlSCglWGHGPNdHvTIiVikL",    "network_wifi": "fEVkCbwiUWeQKBxKUaituWKOb",    "screen_width": "nmlPBglFutdJYDbICgIWhfytF"}



*/

// Reports struct is a row record of the reports table in the tca database
type Reports struct {
	//[ 0] report                                         int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Report int32 `gorm:"primary_key;AUTO_INCREMENT;column:report;type:int;" json:"report"`
	//[ 1] device                                         int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Device null.Int `gorm:"column:device;type:int;" json:"device"`
	//[ 2] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 3] fixed                                          char(5)              null: false  primary: false  isArray: false  auto: false  col: char            len: 5       default: [false]
	Fixed string `gorm:"column:fixed;type:char;size:5;default:false;" json:"fixed"`
	//[ 4] issue                                          int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Issue null.Int `gorm:"column:issue;type:int;" json:"issue"`
	//[ 5] stacktrace                                     text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Stacktrace string `gorm:"column:stacktrace;type:text;size:16777215;" json:"stacktrace"`
	//[ 6] stacktraceGroup                                text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	StacktraceGroup null.String `gorm:"column:stacktraceGroup;type:text;size:65535;" json:"stacktrace_group"`
	//[ 7] log                                            text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Log string `gorm:"column:log;type:text;size:16777215;" json:"log"`
	//[ 8] package                                        text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Package string `gorm:"column:package;type:text;size:16777215;" json:"package"`
	//[ 9] packageVersion                                 text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	PackageVersion string `gorm:"column:packageVersion;type:text;size:16777215;" json:"package_version"`
	//[10] packageVersionCode                             int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [-1]
	PackageVersionCode int32 `gorm:"column:packageVersionCode;type:int;default:-1;" json:"package_version_code"`
	//[11] model                                          text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Model string `gorm:"column:model;type:text;size:16777215;" json:"model"`
	//[12] osVersion                                      text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	OsVersion string `gorm:"column:osVersion;type:text;size:16777215;" json:"os_version"`
	//[13] networkWifi                                    varchar(10)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 10      default: []
	NetworkWifi string `gorm:"column:networkWifi;type:varchar;size:10;" json:"network_wifi"`
	//[14] networkMobile                                  varchar(10)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 10      default: []
	NetworkMobile string `gorm:"column:networkMobile;type:varchar;size:10;" json:"network_mobile"`
	//[15] gps                                            varchar(10)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 10      default: []
	Gps string `gorm:"column:gps;type:varchar;size:10;" json:"gps"`
	//[16] screenWidth                                    varchar(100)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	ScreenWidth string `gorm:"column:screenWidth;type:varchar;size:100;" json:"screen_width"`
	//[17] screenHeight                                   varchar(100)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	ScreenHeight string `gorm:"column:screenHeight;type:varchar;size:100;" json:"screen_height"`
	//[18] screenOrientation                              varchar(100)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	ScreenOrientation string `gorm:"column:screenOrientation;type:varchar;size:100;" json:"screen_orientation"`
	//[19] screenDpi                                      varchar(100)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	ScreenDpi string `gorm:"column:screenDpi;type:varchar;size:100;" json:"screen_dpi"`
}

var reportsTableInfo = &TableInfo{
	Name: "reports",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "report",
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
			GoFieldName:        "Report",
			GoFieldType:        "int32",
			JSONFieldName:      "report",
			ProtobufFieldName:  "report",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "device",
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
			GoFieldName:        "Device",
			GoFieldType:        "null.Int",
			JSONFieldName:      "device",
			ProtobufFieldName:  "device",
			ProtobufType:       "int32",
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
			Name:               "fixed",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "char",
			DatabaseTypePretty: "char(5)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "char",
			ColumnLength:       5,
			GoFieldName:        "Fixed",
			GoFieldType:        "string",
			JSONFieldName:      "fixed",
			ProtobufFieldName:  "fixed",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "issue",
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
			GoFieldName:        "Issue",
			GoFieldType:        "null.Int",
			JSONFieldName:      "issue",
			ProtobufFieldName:  "issue",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "stacktrace",
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
			GoFieldName:        "Stacktrace",
			GoFieldType:        "string",
			JSONFieldName:      "stacktrace",
			ProtobufFieldName:  "stacktrace",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "stacktraceGroup",
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
			GoFieldName:        "StacktraceGroup",
			GoFieldType:        "null.String",
			JSONFieldName:      "stacktrace_group",
			ProtobufFieldName:  "stacktrace_group",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "log",
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
			GoFieldName:        "Log",
			GoFieldType:        "string",
			JSONFieldName:      "log",
			ProtobufFieldName:  "log",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "package",
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
			GoFieldName:        "Package",
			GoFieldType:        "string",
			JSONFieldName:      "package",
			ProtobufFieldName:  "package",
			ProtobufType:       "string",
			ProtobufPos:        9,
		},

		&ColumnInfo{
			Index:              9,
			Name:               "packageVersion",
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
			GoFieldName:        "PackageVersion",
			GoFieldType:        "string",
			JSONFieldName:      "package_version",
			ProtobufFieldName:  "package_version",
			ProtobufType:       "string",
			ProtobufPos:        10,
		},

		&ColumnInfo{
			Index:              10,
			Name:               "packageVersionCode",
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
			GoFieldName:        "PackageVersionCode",
			GoFieldType:        "int32",
			JSONFieldName:      "package_version_code",
			ProtobufFieldName:  "package_version_code",
			ProtobufType:       "int32",
			ProtobufPos:        11,
		},

		&ColumnInfo{
			Index:              11,
			Name:               "model",
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
			GoFieldName:        "Model",
			GoFieldType:        "string",
			JSONFieldName:      "model",
			ProtobufFieldName:  "model",
			ProtobufType:       "string",
			ProtobufPos:        12,
		},

		&ColumnInfo{
			Index:              12,
			Name:               "osVersion",
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
			GoFieldName:        "OsVersion",
			GoFieldType:        "string",
			JSONFieldName:      "os_version",
			ProtobufFieldName:  "os_version",
			ProtobufType:       "string",
			ProtobufPos:        13,
		},

		&ColumnInfo{
			Index:              13,
			Name:               "networkWifi",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(10)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       10,
			GoFieldName:        "NetworkWifi",
			GoFieldType:        "string",
			JSONFieldName:      "network_wifi",
			ProtobufFieldName:  "network_wifi",
			ProtobufType:       "string",
			ProtobufPos:        14,
		},

		&ColumnInfo{
			Index:              14,
			Name:               "networkMobile",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(10)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       10,
			GoFieldName:        "NetworkMobile",
			GoFieldType:        "string",
			JSONFieldName:      "network_mobile",
			ProtobufFieldName:  "network_mobile",
			ProtobufType:       "string",
			ProtobufPos:        15,
		},

		&ColumnInfo{
			Index:              15,
			Name:               "gps",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(10)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       10,
			GoFieldName:        "Gps",
			GoFieldType:        "string",
			JSONFieldName:      "gps",
			ProtobufFieldName:  "gps",
			ProtobufType:       "string",
			ProtobufPos:        16,
		},

		&ColumnInfo{
			Index:              16,
			Name:               "screenWidth",
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
			GoFieldName:        "ScreenWidth",
			GoFieldType:        "string",
			JSONFieldName:      "screen_width",
			ProtobufFieldName:  "screen_width",
			ProtobufType:       "string",
			ProtobufPos:        17,
		},

		&ColumnInfo{
			Index:              17,
			Name:               "screenHeight",
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
			GoFieldName:        "ScreenHeight",
			GoFieldType:        "string",
			JSONFieldName:      "screen_height",
			ProtobufFieldName:  "screen_height",
			ProtobufType:       "string",
			ProtobufPos:        18,
		},

		&ColumnInfo{
			Index:              18,
			Name:               "screenOrientation",
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
			GoFieldName:        "ScreenOrientation",
			GoFieldType:        "string",
			JSONFieldName:      "screen_orientation",
			ProtobufFieldName:  "screen_orientation",
			ProtobufType:       "string",
			ProtobufPos:        19,
		},

		&ColumnInfo{
			Index:              19,
			Name:               "screenDpi",
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
			GoFieldName:        "ScreenDpi",
			GoFieldType:        "string",
			JSONFieldName:      "screen_dpi",
			ProtobufFieldName:  "screen_dpi",
			ProtobufType:       "string",
			ProtobufPos:        20,
		},
	},
}

// TableName sets the insert table name for this struct type
func (r *Reports) TableName() string {
	return "reports"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *Reports) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *Reports) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *Reports) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (r *Reports) TableInfo() *TableInfo {
	return reportsTableInfo
}
