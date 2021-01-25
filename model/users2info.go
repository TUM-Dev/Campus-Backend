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


CREATE TABLE `users2info` (
  `user` int NOT NULL,
  `firstname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `surname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `lastPwChange` int NOT NULL,
  `pager` int DEFAULT '15',
  PRIMARY KEY (`user`),
  CONSTRAINT `fkUsers` FOREIGN KEY (`user`) REFERENCES `users` (`user`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "user": 67,    "firstname": "onqQpeHbIGknZnqxcPJEYThAo",    "surname": "CWmuskKlAPvHuhpWRRIMckDDR",    "last_pw_change": 84,    "pager": 66}



*/

// Users2info struct is a row record of the users2info table in the tca database
type Users2info struct {
	//[ 0] user                                           int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	User int32 `gorm:"primary_key;column:user;type:int;" json:"user"`
	//[ 1] firstname                                      varchar(255)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 255     default: []
	Firstname string `gorm:"column:firstname;type:varchar;size:255;" json:"firstname"`
	//[ 2] surname                                        varchar(255)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 255     default: []
	Surname string `gorm:"column:surname;type:varchar;size:255;" json:"surname"`
	//[ 3] lastPwChange                                   int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	LastPwChange int32 `gorm:"column:lastPwChange;type:int;" json:"last_pw_change"`
	//[ 4] pager                                          int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: [15]
	Pager null.Int `gorm:"column:pager;type:int;default:15;" json:"pager"`
}

var users2infoTableInfo = &TableInfo{
	Name: "users2info",
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
			IsAutoIncrement:    false,
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
			Name:               "firstname",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       255,
			GoFieldName:        "Firstname",
			GoFieldType:        "string",
			JSONFieldName:      "firstname",
			ProtobufFieldName:  "firstname",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "surname",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(255)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       255,
			GoFieldName:        "Surname",
			GoFieldType:        "string",
			JSONFieldName:      "surname",
			ProtobufFieldName:  "surname",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "lastPwChange",
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
			GoFieldName:        "LastPwChange",
			GoFieldType:        "int32",
			JSONFieldName:      "last_pw_change",
			ProtobufFieldName:  "last_pw_change",
			ProtobufType:       "int32",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "pager",
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
			GoFieldName:        "Pager",
			GoFieldType:        "null.Int",
			JSONFieldName:      "pager",
			ProtobufFieldName:  "pager",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (u *Users2info) TableName() string {
	return "users2info"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (u *Users2info) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (u *Users2info) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (u *Users2info) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (u *Users2info) TableInfo() *TableInfo {
	return users2infoTableInfo
}
