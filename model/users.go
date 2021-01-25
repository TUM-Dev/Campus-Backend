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


CREATE TABLE `users` (
  `user` int NOT NULL AUTO_INCREMENT,
  `username` varchar(7) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `firstname` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `surname` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted` int NOT NULL DEFAULT '0',
  `lastActive` int NOT NULL DEFAULT '0',
  `lastPage` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `lastLogin` datetime DEFAULT NULL,
  `tum_id_student` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'OBFUSCATED_ID_ST',
  `tum_id_employee` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'OBFUSCATED_ID_B',
  `tum_id_alumni` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'OBFUSCATED_ID_EXT',
  `tum_id_preferred` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'OBFUSCATED_ID_BEVORZUGT',
  PRIMARY KEY (`user`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=388 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "tum_id_student": "WHDBJCOKdHBRAXojusGEMbPHE",    "tum_id_alumni": "vePhimPhGbwvlpUlXqjuTDDXS",    "firstname": "tEBibkhAQxrWORfuljAGSupWC",    "last_active": 28,    "surname": "frCrbIregwFDbtNGcgNfStOSk",    "created": "2083-05-19T03:55:57.588578815+01:00",    "deleted": 32,    "last_page": "hwiGSiIAYKSuNVHhAFTvdqjDj",    "last_login": "2121-11-27T21:03:21.104789644+01:00",    "tum_id_employee": "ZBfoiboTXQUqoBmFaSOAcYLuJ",    "user": 74,    "username": "ZqWFRDoamxZMceZFZkXAtwJwm",    "tum_id_preferred": "xwhbdJvEvsncwhjLXHmxhXtjU"}



*/

// Users struct is a row record of the users table in the tca database
type Users struct {
	//[ 0] user                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	User int32 `gorm:"primary_key;AUTO_INCREMENT;column:user;type:int;" json:"user"`
	//[ 1] username                                       varchar(7)           null: true   primary: false  isArray: false  auto: false  col: varchar         len: 7       default: []
	Username null.String `gorm:"column:username;type:varchar;size:7;" json:"username"`
	//[ 2] firstname                                      varchar(100)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	Firstname null.String `gorm:"column:firstname;type:varchar;size:100;" json:"firstname"`
	//[ 3] surname                                        varchar(100)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	Surname null.String `gorm:"column:surname;type:varchar;size:100;" json:"surname"`
	//[ 4] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 5] deleted                                        int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	Deleted int32 `gorm:"column:deleted;type:int;default:0;" json:"deleted"`
	//[ 6] lastActive                                     int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	LastActive int32 `gorm:"column:lastActive;type:int;default:0;" json:"last_active"`
	//[ 7] lastPage                                       text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	LastPage null.String `gorm:"column:lastPage;type:text;size:65535;" json:"last_page"`
	//[ 8] lastLogin                                      datetime             null: true   primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	LastLogin null.Time `gorm:"column:lastLogin;type:datetime;" json:"last_login"`
	//[ 9] tum_id_student                                 varchar(50)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 50      default: []
	TumIDStudent null.String `gorm:"column:tum_id_student;type:varchar;size:50;" json:"tum_id_student"` // OBFUSCATED_ID_ST
	//[10] tum_id_employee                                varchar(50)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 50      default: []
	TumIDEmployee null.String `gorm:"column:tum_id_employee;type:varchar;size:50;" json:"tum_id_employee"` // OBFUSCATED_ID_B
	//[11] tum_id_alumni                                  varchar(50)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 50      default: []
	TumIDAlumni null.String `gorm:"column:tum_id_alumni;type:varchar;size:50;" json:"tum_id_alumni"` // OBFUSCATED_ID_EXT
	//[12] tum_id_preferred                               varchar(50)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 50      default: []
	TumIDPreferred null.String `gorm:"column:tum_id_preferred;type:varchar;size:50;" json:"tum_id_preferred"` // OBFUSCATED_ID_BEVORZUGT

}

var usersTableInfo = &TableInfo{
	Name: "users",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "user",
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
			GoFieldName:        "User",
			GoFieldType:        "int32",
			JSONFieldName:      "user",
			ProtobufFieldName:  "user",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "username",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(7)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       7,
			GoFieldName:        "Username",
			GoFieldType:        "null.String",
			JSONFieldName:      "username",
			ProtobufFieldName:  "username",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "firstname",
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
			GoFieldName:        "Firstname",
			GoFieldType:        "null.String",
			JSONFieldName:      "firstname",
			ProtobufFieldName:  "firstname",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "surname",
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
			GoFieldName:        "Surname",
			GoFieldType:        "null.String",
			JSONFieldName:      "surname",
			ProtobufFieldName:  "surname",
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
			Name:               "deleted",
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
			GoFieldName:        "Deleted",
			GoFieldType:        "int32",
			JSONFieldName:      "deleted",
			ProtobufFieldName:  "deleted",
			ProtobufType:       "int32",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "lastActive",
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
			GoFieldName:        "LastActive",
			GoFieldType:        "int32",
			JSONFieldName:      "last_active",
			ProtobufFieldName:  "last_active",
			ProtobufType:       "int32",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "lastPage",
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
			GoFieldName:        "LastPage",
			GoFieldType:        "null.String",
			JSONFieldName:      "last_page",
			ProtobufFieldName:  "last_page",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},

		&ColumnInfo{
			Index:              8,
			Name:               "lastLogin",
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
			GoFieldName:        "LastLogin",
			GoFieldType:        "null.Time",
			JSONFieldName:      "last_login",
			ProtobufFieldName:  "last_login",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        9,
		},

		&ColumnInfo{
			Index:              9,
			Name:               "tum_id_student",
			Comment:            `OBFUSCATED_ID_ST`,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(50)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       50,
			GoFieldName:        "TumIDStudent",
			GoFieldType:        "null.String",
			JSONFieldName:      "tum_id_student",
			ProtobufFieldName:  "tum_id_student",
			ProtobufType:       "string",
			ProtobufPos:        10,
		},

		&ColumnInfo{
			Index:              10,
			Name:               "tum_id_employee",
			Comment:            `OBFUSCATED_ID_B`,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(50)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       50,
			GoFieldName:        "TumIDEmployee",
			GoFieldType:        "null.String",
			JSONFieldName:      "tum_id_employee",
			ProtobufFieldName:  "tum_id_employee",
			ProtobufType:       "string",
			ProtobufPos:        11,
		},

		&ColumnInfo{
			Index:              11,
			Name:               "tum_id_alumni",
			Comment:            `OBFUSCATED_ID_EXT`,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(50)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       50,
			GoFieldName:        "TumIDAlumni",
			GoFieldType:        "null.String",
			JSONFieldName:      "tum_id_alumni",
			ProtobufFieldName:  "tum_id_alumni",
			ProtobufType:       "string",
			ProtobufPos:        12,
		},

		&ColumnInfo{
			Index:              12,
			Name:               "tum_id_preferred",
			Comment:            `OBFUSCATED_ID_BEVORZUGT`,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(50)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       50,
			GoFieldName:        "TumIDPreferred",
			GoFieldType:        "null.String",
			JSONFieldName:      "tum_id_preferred",
			ProtobufFieldName:  "tum_id_preferred",
			ProtobufType:       "string",
			ProtobufPos:        13,
		},
	},
}

// TableName sets the insert table name for this struct type
func (u *Users) TableName() string {
	return "users"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (u *Users) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (u *Users) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (u *Users) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (u *Users) TableInfo() *TableInfo {
	return usersTableInfo
}
