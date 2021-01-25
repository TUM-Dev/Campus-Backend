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


CREATE TABLE `barrierFree_moreInfo` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(32) DEFAULT NULL,
  `category` varchar(11) DEFAULT NULL,
  `url` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8

JSON Sample
-------------------------------------
{    "url": "tIKlbJXxfaVBbqiQXMfIFtNqv",    "id": 71,    "title": "qWsURFpmyCyelUCwYxZvdcWJm",    "category": "raiXRHMTmMCThjpSPZIgFYFii"}


Comments
-------------------------------------
[ 0] column is set for unsigned



*/

// BarrierFreeMoreInfo struct is a row record of the barrierFree_moreInfo table in the tca database
type BarrierFreeMoreInfo struct {
	//[ 0] id                                             uint                 null: false  primary: true   isArray: false  auto: true   col: uint            len: -1      default: []
	ID uint32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:uint;" json:"id"`
	//[ 1] title                                          varchar(32)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	Title null.String `gorm:"column:title;type:varchar;size:32;" json:"title"`
	//[ 2] category                                       varchar(11)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 11      default: []
	Category null.String `gorm:"column:category;type:varchar;size:11;" json:"category"`
	//[ 3] url                                            varchar(128)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 128     default: []
	URL null.String `gorm:"column:url;type:varchar;size:128;" json:"url"`
}

var barrierFree_moreInfoTableInfo = &TableInfo{
	Name: "barrierFree_moreInfo",
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
			Name:               "title",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(32)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       32,
			GoFieldName:        "Title",
			GoFieldType:        "null.String",
			JSONFieldName:      "title",
			ProtobufFieldName:  "title",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "category",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(11)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       11,
			GoFieldName:        "Category",
			GoFieldType:        "null.String",
			JSONFieldName:      "category",
			ProtobufFieldName:  "category",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "url",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(128)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       128,
			GoFieldName:        "URL",
			GoFieldType:        "null.String",
			JSONFieldName:      "url",
			ProtobufFieldName:  "url",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},
	},
}

// TableName sets the insert table name for this struct type
func (b *BarrierFreeMoreInfo) TableName() string {
	return "barrierFree_moreInfo"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (b *BarrierFreeMoreInfo) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (b *BarrierFreeMoreInfo) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (b *BarrierFreeMoreInfo) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (b *BarrierFreeMoreInfo) TableInfo() *TableInfo {
	return barrierFree_moreInfoTableInfo
}
