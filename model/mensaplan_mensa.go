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


CREATE TABLE `mensaplan_mensa` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `latitude` double DEFAULT NULL,
  `longitude` double DEFAULT NULL,
  `webid` int DEFAULT NULL,
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=30 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

JSON Sample
-------------------------------------
{    "category": "McdwFEjfAVBKjZMTPpnaLZtKB",    "id": 74,    "name": "qnIncBNQATYJRbroRIFrtagZi",    "latitude": 0.3868265515741528,    "longitude": 0.5604240939820241,    "webid": 10}



*/

// MensaplanMensa struct is a row record of the mensaplan_mensa table in the tca database
type MensaplanMensa struct {
	//[ 0] id                                             int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	ID int32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id"`
	//[ 1] name                                           varchar(100)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	Name string `gorm:"column:name;type:varchar;size:100;" json:"name"`
	//[ 2] latitude                                       double               null: true   primary: false  isArray: false  auto: false  col: double          len: -1      default: []
	Latitude null.Float `gorm:"column:latitude;type:double;" json:"latitude"`
	//[ 3] longitude                                      double               null: true   primary: false  isArray: false  auto: false  col: double          len: -1      default: []
	Longitude null.Float `gorm:"column:longitude;type:double;" json:"longitude"`
	//[ 4] webid                                          int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Webid null.Int `gorm:"column:webid;type:int;" json:"webid"`
	//[ 5] category                                       varchar(50)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 50      default: []
	Category string `gorm:"column:category;type:varchar;size:50;" json:"category"`
}

var mensaplan_mensaTableInfo = &TableInfo{
	Name: "mensaplan_mensa",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
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
			GoFieldName:        "ID",
			GoFieldType:        "int32",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "name",
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
			GoFieldName:        "Name",
			GoFieldType:        "string",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "latitude",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "double",
			DatabaseTypePretty: "double",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "double",
			ColumnLength:       -1,
			GoFieldName:        "Latitude",
			GoFieldType:        "null.Float",
			JSONFieldName:      "latitude",
			ProtobufFieldName:  "latitude",
			ProtobufType:       "float",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "longitude",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "double",
			DatabaseTypePretty: "double",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "double",
			ColumnLength:       -1,
			GoFieldName:        "Longitude",
			GoFieldType:        "null.Float",
			JSONFieldName:      "longitude",
			ProtobufFieldName:  "longitude",
			ProtobufType:       "float",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "webid",
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
			GoFieldName:        "Webid",
			GoFieldType:        "null.Int",
			JSONFieldName:      "webid",
			ProtobufFieldName:  "webid",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "category",
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
			GoFieldName:        "Category",
			GoFieldType:        "string",
			JSONFieldName:      "category",
			ProtobufFieldName:  "category",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},
	},
}

// TableName sets the insert table name for this struct type
func (m *MensaplanMensa) TableName() string {
	return "mensaplan_mensa"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (m *MensaplanMensa) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (m *MensaplanMensa) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (m *MensaplanMensa) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (m *MensaplanMensa) TableInfo() *TableInfo {
	return mensaplan_mensaTableInfo
}
