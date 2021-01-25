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


CREATE TABLE `roomfinder_buildings2gps` (
  `id` varchar(8) NOT NULL DEFAULT '',
  `latitude` varchar(30) DEFAULT NULL,
  `longitude` varchar(30) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

JSON Sample
-------------------------------------
{    "longitude": "yxBCFuUSdyxBiAbYJtHheSkAi",    "id": "rVFqlYjfDoUOShXtwpJGcuYSj",    "latitude": "pIeDIOFBNOukjvGfVWmakSWOt"}



*/

// RoomfinderBuildings2gps struct is a row record of the roomfinder_buildings2gps table in the tca database
type RoomfinderBuildings2gps struct {
	//[ 0] id                                             varchar(8)           null: false  primary: true   isArray: false  auto: false  col: varchar         len: 8       default: []
	ID string `gorm:"primary_key;column:id;type:varchar;size:8;" json:"id"`
	//[ 1] latitude                                       varchar(30)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 30      default: []
	Latitude null.String `gorm:"column:latitude;type:varchar;size:30;" json:"latitude"`
	//[ 2] longitude                                      varchar(30)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 30      default: []
	Longitude null.String `gorm:"column:longitude;type:varchar;size:30;" json:"longitude"`
}

var roomfinder_buildings2gpsTableInfo = &TableInfo{
	Name: "roomfinder_buildings2gps",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(8)",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       8,
			GoFieldName:        "ID",
			GoFieldType:        "string",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "latitude",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(30)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       30,
			GoFieldName:        "Latitude",
			GoFieldType:        "null.String",
			JSONFieldName:      "latitude",
			ProtobufFieldName:  "latitude",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "longitude",
			Comment:            ``,
			Notes:              ``,
			Nullable:           true,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(30)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       30,
			GoFieldName:        "Longitude",
			GoFieldType:        "null.String",
			JSONFieldName:      "longitude",
			ProtobufFieldName:  "longitude",
			ProtobufType:       "string",
			ProtobufPos:        3,
		},
	},
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderBuildings2gps) TableName() string {
	return "roomfinder_buildings2gps"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *RoomfinderBuildings2gps) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *RoomfinderBuildings2gps) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *RoomfinderBuildings2gps) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (r *RoomfinderBuildings2gps) TableInfo() *TableInfo {
	return roomfinder_buildings2gpsTableInfo
}
