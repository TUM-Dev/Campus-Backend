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


CREATE TABLE `barrierFree_persons` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(40) DEFAULT NULL,
  `telephone` varchar(32) DEFAULT NULL,
  `email` varchar(32) DEFAULT NULL,
  `faculty` varchar(32) DEFAULT NULL,
  `office` varchar(16) DEFAULT NULL,
  `officeHour` varchar(16) DEFAULT NULL,
  `tumID` varchar(24) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8

JSON Sample
-------------------------------------
{    "faculty": "xcIGeNMNlxqMiHACrGjUEYXtm",    "office": "gqQnamFaYwoLKyaoeexlHkvwI",    "office_hour": "jqFZiHsXtGwEIInhNCrLLGMqd",    "tum_id": "uxpIVuaekeussZHUdAErSwYXR",    "id": 56,    "name": "ujfEulTbDKTJTaBdnxNqAbGhf",    "telephone": "IAfnwoDtYYRePHyENQSHUZoin",    "email": "JCwZdtsnodRKFWGEVhOfLalLh"}


Comments
-------------------------------------
[ 0] column is set for unsigned



*/

// BarrierFreePersons struct is a row record of the barrierFree_persons table in the tca database
type BarrierFreePersons struct {
	//[ 0] id                                             uint                 null: false  primary: true   isArray: false  auto: true   col: uint            len: -1      default: []
	ID uint32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:uint;" json:"id"`
	//[ 1] name                                           varchar(40)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 40      default: []
	Name null.String `gorm:"column:name;type:varchar;size:40;" json:"name"`
	//[ 2] telephone                                      varchar(32)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	Telephone null.String `gorm:"column:telephone;type:varchar;size:32;" json:"telephone"`
	//[ 3] email                                          varchar(32)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	Email null.String `gorm:"column:email;type:varchar;size:32;" json:"email"`
	//[ 4] faculty                                        varchar(32)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	Faculty null.String `gorm:"column:faculty;type:varchar;size:32;" json:"faculty"`
	//[ 5] office                                         varchar(16)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 16      default: []
	Office null.String `gorm:"column:office;type:varchar;size:16;" json:"office"`
	//[ 6] officeHour                                     varchar(16)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 16      default: []
	OfficeHour null.String `gorm:"column:officeHour;type:varchar;size:16;" json:"office_hour"`
	//[ 7] tumID                                          varchar(24)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 24      default: []
	TumID null.String `gorm:"column:tumID;type:varchar;size:24;" json:"tum_id"`
}

var barrierFree_personsTableInfo = &TableInfo{
	Name: "barrierFree_persons",
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
			Name:               "name",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(40)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       40,
			GoFieldName:        "Name",
			GoFieldType:        "null.String",
			JSONFieldName:      "name",
			ProtobufFieldName:  "name",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "telephone",
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
			GoFieldName:        "Telephone",
			GoFieldType:        "null.String",
			JSONFieldName:      "telephone",
			ProtobufFieldName:  "telephone",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "email",
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
			GoFieldName:        "Email",
			GoFieldType:        "null.String",
			JSONFieldName:      "email",
			ProtobufFieldName:  "email",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "faculty",
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
			GoFieldName:        "Faculty",
			GoFieldType:        "null.String",
			JSONFieldName:      "faculty",
			ProtobufFieldName:  "faculty",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "office",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(16)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       16,
			GoFieldName:        "Office",
			GoFieldType:        "null.String",
			JSONFieldName:      "office",
			ProtobufFieldName:  "office",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "officeHour",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(16)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       16,
			GoFieldName:        "OfficeHour",
			GoFieldType:        "null.String",
			JSONFieldName:      "office_hour",
			ProtobufFieldName:  "office_hour",
			ProtobufType:       "string",
			ProtobufPos:        7,
		},

		&ColumnInfo{
			Index:              7,
			Name:               "tumID",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(24)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       24,
			GoFieldName:        "TumID",
			GoFieldType:        "null.String",
			JSONFieldName:      "tum_id",
			ProtobufFieldName:  "tum_id",
			ProtobufType:       "string",
			ProtobufPos:        8,
		},
	},
}

// TableName sets the insert table name for this struct type
func (b *BarrierFreePersons) TableName() string {
	return "barrierFree_persons"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (b *BarrierFreePersons) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (b *BarrierFreePersons) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (b *BarrierFreePersons) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (b *BarrierFreePersons) TableInfo() *TableInfo {
	return barrierFree_personsTableInfo
}
