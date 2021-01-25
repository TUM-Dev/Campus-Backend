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


CREATE TABLE `recover` (
  `recover` int NOT NULL AUTO_INCREMENT,
  `user` int NOT NULL,
  `created` int NOT NULL,
  `hash` varchar(190) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `ip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`recover`),
  UNIQUE KEY `hash` (`hash`),
  KEY `user` (`user`),
  CONSTRAINT `fkRecover2User` FOREIGN KEY (`user`) REFERENCES `users` (`user`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "recover": 15,    "user": 66,    "created": 74,    "hash": "lESDEVXFEWZnDfGHLcpiPCMGF",    "ip": "BCgYnLDTgbSKkOVZnlQwcoLJh"}



*/

// Recover struct is a row record of the recover table in the tca database
type Recover struct {
	//[ 0] recover                                        int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Recover int32 `gorm:"primary_key;AUTO_INCREMENT;column:recover;type:int;" json:"recover"`
	//[ 1] user                                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	User int32 `gorm:"column:user;type:int;" json:"user"`
	//[ 2] created                                        int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Created int32 `gorm:"column:created;type:int;" json:"created"`
	//[ 3] hash                                           varchar(190)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 190     default: []
	Hash string `gorm:"column:hash;type:varchar;size:190;" json:"hash"`
	//[ 4] ip                                             varchar(255)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 255     default: []
	IP string `gorm:"column:ip;type:varchar;size:255;" json:"ip"`
}

var recoverTableInfo = &TableInfo{
	Name: "recover",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "recover",
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
			GoFieldName:        "Recover",
			GoFieldType:        "int32",
			JSONFieldName:      "recover",
			ProtobufFieldName:  "recover",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "user",
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
			GoFieldName:        "User",
			GoFieldType:        "int32",
			JSONFieldName:      "user",
			ProtobufFieldName:  "user",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "created",
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
			GoFieldName:        "Created",
			GoFieldType:        "int32",
			JSONFieldName:      "created",
			ProtobufFieldName:  "created",
			ProtobufType:       "int32",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "hash",
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
			GoFieldName:        "Hash",
			GoFieldType:        "string",
			JSONFieldName:      "hash",
			ProtobufFieldName:  "hash",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "ip",
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
			GoFieldName:        "IP",
			GoFieldType:        "string",
			JSONFieldName:      "ip",
			ProtobufFieldName:  "ip",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},
	},
}

// TableName sets the insert table name for this struct type
func (r *Recover) TableName() string {
	return "recover"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *Recover) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *Recover) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *Recover) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (r *Recover) TableInfo() *TableInfo {
	return recoverTableInfo
}
