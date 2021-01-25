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


CREATE TABLE `member` (
  `member` int NOT NULL AUTO_INCREMENT,
  `lrz_id` varchar(7) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `active_day` int DEFAULT '0',
  `active_day_date` date DEFAULT NULL,
  `student_id` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `employee_id` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `external_id` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`member`),
  UNIQUE KEY `lrz_id` (`lrz_id`)
) ENGINE=InnoDB AUTO_INCREMENT=92496 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "lrz_id": "pBBIlSWbJLiJxLsvjtULMMItl",    "name": "EwgALmQmQJkNJXwWjcffMTHoi",    "active_day": 69,    "active_day_date": "2289-07-15T01:57:07.733162553+01:00",    "student_id": "HmhwAXkvMaGWQgfwJOFMCdYWU",    "employee_id": "lwEgjLPKxXLXtvVgpVpCNFbSn",    "external_id": "xRoHWTZujIlwaIwunnoGbSxuh",    "member": 26}



*/

// Member struct is a row record of the member table in the tca database
type Member struct {
	//[ 0] member                                         int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Member int32 `gorm:"primary_key;AUTO_INCREMENT;column:member;type:int;" json:"member"`
	//[ 1] lrz_id                                         varchar(7)           null: false  primary: false  isArray: false  auto: false  col: varchar         len: 7       default: []
	LrzID string `gorm:"column:lrz_id;type:varchar;size:7;" json:"lrz_id"`
	//[ 2] name                                           varchar(150)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 150     default: []
	Name string `gorm:"column:name;type:varchar;size:150;" json:"name"`
	//[ 3] active_day                                     int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	ActiveDay null.Int `gorm:"column:active_day;type:int;default:0;" json:"active_day"`
	//[ 4] active_day_date                                date                 null: true   primary: false  isArray: false  auto: false  col: date            len: -1      default: []
	ActiveDayDate null.Time `gorm:"column:active_day_date;type:date;" json:"active_day_date"`
	//[ 5] student_id                                     text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	StudentID null.String `gorm:"column:student_id;type:text;size:65535;" json:"student_id"`
	//[ 6] employee_id                                    text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	EmployeeID null.String `gorm:"column:employee_id;type:text;size:65535;" json:"employee_id"`
	//[ 7] external_id                                    text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	ExternalID null.String `gorm:"column:external_id;type:text;size:65535;" json:"external_id"`
}

var memberTableInfo = &TableInfo{
	Name: "member",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "member",
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
			GoFieldName:        "Member",
			GoFieldType:        "int32",
			JSONFieldName:      "member",
			ProtobufFieldName:  "member",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "lrz_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(7)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       7,
			GoFieldName:        "LrzID",
			GoFieldType:        "string",
			JSONFieldName:      "lrz_id",
			ProtobufFieldName:  "lrz_id",
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
			DatabaseTypePretty: "varchar(150)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       150,
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "active_day",
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
			GoFieldName:        "ActiveDay",
			GoFieldType:        "null.Int",
			JSONFieldName:      "active_day",
			ProtobufFieldName:  "active_day",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "active_day_date",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "date",
			DatabaseTypePretty: "date",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "date",
			ColumnLength:       -1,
			GoFieldName:        "ActiveDayDate",
			GoFieldType:        "null.Time",
			JSONFieldName:      "active_day_date",
			ProtobufFieldName:  "active_day_date",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "student_id",
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
			GoFieldName:        "StudentID",
			GoFieldType:        "null.String",
			JSONFieldName:      "student_id",
			ProtobufFieldName:  "student_id",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "employee_id",
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
			GoFieldName:        "EmployeeID",
			GoFieldType:        "null.String",
			JSONFieldName:      "employee_id",
			ProtobufFieldName:  "employee_id",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "external_id",
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
			GoFieldName:        "ExternalID",
			GoFieldType:        "null.String",
			JSONFieldName:      "external_id",
			ProtobufFieldName:  "external_id",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},
	},
}

// TableName sets the insert table name for this struct type
func (m *Member) TableName() string {
	return "member"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (m *Member) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (m *Member) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (m *Member) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (m *Member) TableInfo() *TableInfo {
	return memberTableInfo
}
