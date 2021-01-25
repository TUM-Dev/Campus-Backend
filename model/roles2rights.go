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


CREATE TABLE `roles2rights` (
  `role` int NOT NULL,
  `right` int NOT NULL,
  PRIMARY KEY (`role`,`right`),
  KEY `fkRight_idx` (`right`),
  CONSTRAINT `fkRight` FOREIGN KEY (`right`) REFERENCES `rights` (`right`) ON DELETE CASCADE,
  CONSTRAINT `fkRole` FOREIGN KEY (`role`) REFERENCES `roles` (`role`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "role": 61,    "right": 78}



*/

// Roles2rights struct is a row record of the roles2rights table in the tca database
type Roles2rights struct {
	//[ 0] role                                           int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	Role int32 `gorm:"primary_key;column:role;type:int;" json:"role"`
	//[ 1] right                                          int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	Right int32 `gorm:"primary_key;column:right;type:int;" json:"right"`
}

var roles2rightsTableInfo = &TableInfo{
	Name: "roles2rights",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "role",
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
			GoFieldName:        "Role",
			GoFieldType:        "int32",
			JSONFieldName:      "role",
			ProtobufFieldName:  "role",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "right",
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
			GoFieldName:        "Right",
			GoFieldType:        "int32",
			JSONFieldName:      "right",
			ProtobufFieldName:  "right",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (r *Roles2rights) TableName() string {
	return "roles2rights"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *Roles2rights) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *Roles2rights) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *Roles2rights) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (r *Roles2rights) TableInfo() *TableInfo {
	return roles2rightsTableInfo
}
