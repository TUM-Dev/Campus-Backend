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


CREATE TABLE `dishflags` (
  `flag` int NOT NULL AUTO_INCREMENT,
  `short` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`flag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "flag": 0,    "short": "dBXbcjKUFouiitYNCEneteiTF",    "description": "axFcpWvOlhVIEPbeUjWnInLrE"}



*/

// Dishflags struct is a row record of the dishflags table in the tca database
type Dishflags struct {
	//[ 0] flag                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Flag int32 `gorm:"primary_key;AUTO_INCREMENT;column:flag;type:int;" json:"flag"`
	//[ 1] short                                          varchar(10)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 10      default: []
	Short string `gorm:"column:short;type:varchar;size:10;" json:"short"`
	//[ 2] description                                    varchar(50)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 50      default: []
	Description string `gorm:"column:description;type:varchar;size:50;" json:"description"`
}

var dishflagsTableInfo = &TableInfo{
	Name: "dishflags",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "flag",
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
			GoFieldName:        "Flag",
			GoFieldType:        "int32",
			JSONFieldName:      "flag",
			ProtobufFieldName:  "flag",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "short",
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
			GoFieldName:        "Short",
			GoFieldType:        "string",
			JSONFieldName:      "short",
			ProtobufFieldName:  "short",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "description",
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
			GoFieldName:        "Description",
			GoFieldType:        "string",
			JSONFieldName:      "description",
			ProtobufFieldName:  "description",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (d *Dishflags) TableName() string {
	return "dishflags"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (d *Dishflags) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (d *Dishflags) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (d *Dishflags) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (d *Dishflags) TableInfo() *TableInfo {
	return dishflagsTableInfo
}
