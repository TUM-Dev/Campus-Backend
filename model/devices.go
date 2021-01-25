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


CREATE TABLE `devices` (
  `device` int NOT NULL AUTO_INCREMENT,
  `member` int DEFAULT NULL,
  `uuid` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created` timestamp NULL DEFAULT NULL,
  `lastAccess` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE CURRENT_TIMESTAMP,
  `lastApi` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `developer` enum('true','false') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'false',
  `osVersion` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `appVersion` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `counter` int NOT NULL DEFAULT '0',
  `pk` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `pkActive` enum('true','false') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'false',
  `gcmToken` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `gcmStatus` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `confirmationKey` varchar(35) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `keyCreated` datetime DEFAULT NULL,
  `keyConfirmed` datetime DEFAULT NULL,
  PRIMARY KEY (`device`),
  UNIQUE KEY `uuid` (`uuid`),
  KEY `member` (`member`),
  CONSTRAINT `devices_ibfk_1` FOREIGN KEY (`member`) REFERENCES `member` (`member`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=123932 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "last_api": "pKYlruBKSuZNsAGACHteOFKba",    "pk_active": "PEjXmHYHlTcPnhWTcNJVTsNQE",    "app_version": "LvgNYSiJKMKWUiIwqawAiEFtp",    "device": 87,    "member": 83,    "developer": "tyWXhofALnAlwANjcnkchBLcG",    "pk": "gQqtBRTkqSLliEeuWjvBRIkPT",    "confirmation_key": "VPkMwmNilnjLtBqtUDAeywqwh",    "uuid": "kqntKNsmmZSnJqsTtwlXBoEwg",    "created": "2182-09-05T01:47:00.223911811+01:00",    "counter": 53,    "gcm_status": "maEglsvFsaoLYhCMRNsmbZHCK",    "key_created": "2295-06-02T19:34:46.874128651+01:00",    "key_confirmed": "2299-10-13T17:15:01.137476144+01:00",    "last_access": "2031-05-27T11:05:38.663938565+02:00",    "os_version": "BrpIgCOXORUsbpgGwbKfkhhgQ",    "gcm_token": "OGqHjaWCHUmnGuGfBsqdXnQHi"}



*/

// Devices struct is a row record of the devices table in the tca database
type Devices struct {
	//[ 0] device                                         int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Device int32 `gorm:"primary_key;AUTO_INCREMENT;column:device;type:int;" json:"device"`
	//[ 1] member                                         int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Member null.Int `gorm:"column:member;type:int;" json:"member"`
	//[ 2] uuid                                           varchar(50)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 50      default: []
	UUID string `gorm:"column:uuid;type:varchar;size:50;" json:"uuid"`
	//[ 3] created                                        timestamp            null: true   primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: []
	Created null.Time `gorm:"column:created;type:timestamp;" json:"created"`
	//[ 4] lastAccess                                     timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [0000-00-00 00:00:00]
	LastAccess time.Time `gorm:"column:lastAccess;type:timestamp;default:0000-00-00 00:00:00;" json:"last_access"`
	//[ 5] lastApi                                        text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	LastAPI string `gorm:"column:lastApi;type:text;size:16777215;" json:"last_api"`
	//[ 6] developer                                      char(5)              null: false  primary: false  isArray: false  auto: false  col: char            len: 5       default: [false]
	Developer string `gorm:"column:developer;type:char;size:5;default:false;" json:"developer"`
	//[ 7] osVersion                                      text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	OsVersion string `gorm:"column:osVersion;type:text;size:16777215;" json:"os_version"`
	//[ 8] appVersion                                     text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	AppVersion string `gorm:"column:appVersion;type:text;size:16777215;" json:"app_version"`
	//[ 9] counter                                        int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	Counter int32 `gorm:"column:counter;type:int;default:0;" json:"counter"`
	//[10] pk                                             text(4294967295)     null: true   primary: false  isArray: false  auto: false  col: text            len: 4294967295 default: []
	Pk null.String `gorm:"column:pk;type:text;size:4294967295;" json:"pk"`
	//[11] pkActive                                       char(5)              null: false  primary: false  isArray: false  auto: false  col: char            len: 5       default: [false]
	PkActive string `gorm:"column:pkActive;type:char;size:5;default:false;" json:"pk_active"`
	//[12] gcmToken                                       text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	GcmToken null.String `gorm:"column:gcmToken;type:text;size:65535;" json:"gcm_token"`
	//[13] gcmStatus                                      varchar(200)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 200     default: []
	GcmStatus null.String `gorm:"column:gcmStatus;type:varchar;size:200;" json:"gcm_status"`
	//[14] confirmationKey                                varchar(35)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 35      default: []
	ConfirmationKey null.String `gorm:"column:confirmationKey;type:varchar;size:35;" json:"confirmation_key"`
	//[15] keyCreated                                     datetime             null: true   primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	KeyCreated null.Time `gorm:"column:keyCreated;type:datetime;" json:"key_created"`
	//[16] keyConfirmed                                   datetime             null: true   primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	KeyConfirmed null.Time `gorm:"column:keyConfirmed;type:datetime;" json:"key_confirmed"`
}

var devicesTableInfo = &TableInfo{
	Name: "devices",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "device",
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
			GoFieldName:        "Device",
			GoFieldType:        "int32",
			JSONFieldName:      "device",
			ProtobufFieldName:  "device",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "member",
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
			GoFieldName:        "Member",
			GoFieldType:        "null.Int",
			JSONFieldName:      "member",
			ProtobufFieldName:  "member",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "uuid",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(50)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       50,
			GoFieldName:        "UUID",
			GoFieldType:        "string",
			JSONFieldName:      "uuid",
			ProtobufFieldName:  "uuid",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "created",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "timestamp",
			DatabaseTypePretty: "timestamp",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "timestamp",
			ColumnLength:       -1,
			GoFieldName:        "Created",
			GoFieldType:        "null.Time",
			JSONFieldName:      "created",
			ProtobufFieldName:  "created",
			ProtobufType:       "uint64",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "lastAccess",
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
			GoFieldName:        "LastAccess",
			GoFieldType:        "time.Time",
			JSONFieldName:      "last_access",
			ProtobufFieldName:  "last_access",
			ProtobufType:       "uint64",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "lastApi",
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
			GoFieldName:        "LastAPI",
			GoFieldType:        "string",
			JSONFieldName:      "last_api",
			ProtobufFieldName:  "last_api",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "developer",
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
			GoFieldName:        "Developer",
			GoFieldType:        "string",
			JSONFieldName:      "developer",
			ProtobufFieldName:  "developer",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
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
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "appVersion",
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
			GoFieldName:        "AppVersion",
			GoFieldType:        "string",
			JSONFieldName:      "app_version",
			ProtobufFieldName:  "app_version",
			ProtobufType:       "string",
			ProtobufPos:        9,
		},

		&ColumnInfo{
			Index:              9,
			Name:               "counter",
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
			GoFieldName:        "Counter",
			GoFieldType:        "int32",
			JSONFieldName:      "counter",
			ProtobufFieldName:  "counter",
			ProtobufType:       "int32",
			ProtobufPos:        10,
		},

		&ColumnInfo{
			Index:              10,
			Name:               "pk",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "text",
			DatabaseTypePretty: "text(4294967295)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "text",
			ColumnLength:       4294967295,
			GoFieldName:        "Pk",
			GoFieldType:        "null.String",
			JSONFieldName:      "pk",
			ProtobufFieldName:  "pk",
			ProtobufType:       "string",
			ProtobufPos:        11,
		},

		&ColumnInfo{
			Index:              11,
			Name:               "pkActive",
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
			GoFieldName:        "PkActive",
			GoFieldType:        "string",
			JSONFieldName:      "pk_active",
			ProtobufFieldName:  "pk_active",
			ProtobufType:       "string",
			ProtobufPos:        12,
		},

		&ColumnInfo{
			Index:              12,
			Name:               "gcmToken",
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
			GoFieldName:        "GcmToken",
			GoFieldType:        "null.String",
			JSONFieldName:      "gcm_token",
			ProtobufFieldName:  "gcm_token",
			ProtobufType:       "string",
			ProtobufPos:        13,
		},

		&ColumnInfo{
			Index:              13,
			Name:               "gcmStatus",
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
			GoFieldName:        "GcmStatus",
			GoFieldType:        "null.String",
			JSONFieldName:      "gcm_status",
			ProtobufFieldName:  "gcm_status",
			ProtobufType:       "string",
			ProtobufPos:        14,
		},

		&ColumnInfo{
			Index:              14,
			Name:               "confirmationKey",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(35)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       35,
			GoFieldName:        "ConfirmationKey",
			GoFieldType:        "null.String",
			JSONFieldName:      "confirmation_key",
			ProtobufFieldName:  "confirmation_key",
			ProtobufType:       "string",
			ProtobufPos:        15,
		},

		&ColumnInfo{
			Index:              15,
			Name:               "keyCreated",
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
			GoFieldName:        "KeyCreated",
			GoFieldType:        "null.Time",
			JSONFieldName:      "key_created",
			ProtobufFieldName:  "key_created",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        16,
		},

		&ColumnInfo{
			Index:              16,
			Name:               "keyConfirmed",
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
			GoFieldName:        "KeyConfirmed",
			GoFieldType:        "null.Time",
			JSONFieldName:      "key_confirmed",
			ProtobufFieldName:  "key_confirmed",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        17,
		},
	},
}

// TableName sets the insert table name for this struct type
func (d *Devices) TableName() string {
	return "devices"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *Devices) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *Devices) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *Devices) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *Devices) TableInfo() *TableInfo {
	return devicesTableInfo
}
