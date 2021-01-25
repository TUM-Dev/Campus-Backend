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


CREATE TABLE `roomfinder_schedules` (
  `room_id` int NOT NULL,
  `start` datetime NOT NULL,
  `end` datetime NOT NULL,
  `title` varchar(64) NOT NULL,
  `event_id` int NOT NULL,
  `course_code` varchar(32) DEFAULT NULL,
  UNIQUE KEY `unique` (`room_id`,`start`,`end`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

JSON Sample
-------------------------------------
{    "course_code": "gIryKIKIWMILfOjPmesHSWuml",    "room_id": 22,    "start": "2030-05-06T15:46:46.307587013+02:00",    "end": "2039-02-12T16:35:48.121077202+01:00",    "title": "DxXPVROhiIZdeutJSTrrhstBg",    "event_id": 32}


Comments
-------------------------------------
[ 0] Warning table: roomfinder_schedules does not have a primary key defined, setting col position 1 room_id as primary key




*/

// RoomfinderSchedules struct is a row record of the roomfinder_schedules table in the tca database
type RoomfinderSchedules struct {
	//[ 0] room_id                                        int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	RoomID int32 `gorm:"primary_key;column:room_id;type:int;" json:"room_id"`
	//[ 1] start                                          datetime             null: false  primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	Start time.Time `gorm:"column:start;type:datetime;" json:"start"`
	//[ 2] end                                            datetime             null: false  primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	End time.Time `gorm:"column:end;type:datetime;" json:"end"`
	//[ 3] title                                          varchar(64)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 64      default: []
	Title string `gorm:"column:title;type:varchar;size:64;" json:"title"`
	//[ 4] event_id                                       int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	EventID int32 `gorm:"column:event_id;type:int;" json:"event_id"`
	//[ 5] course_code                                    varchar(32)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 32      default: []
	CourseCode null.String `gorm:"column:course_code;type:varchar;size:32;" json:"course_code"`
}

var roomfinder_schedulesTableInfo = &TableInfo{
	Name: "roomfinder_schedules",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:   0,
			Name:    "room_id",
			Comment: ``,
			Notes: `Warning table: roomfinder_schedules does not have a primary key defined, setting col position 1 room_id as primary key
`,
			Nullable:           false,
			DatabaseTypeName:   "int",
			DatabaseTypePretty: "int",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "int",
			ColumnLength:       -1,
			GoFieldName:        "RoomID",
			GoFieldType:        "int32",
			JSONFieldName:      "room_id",
			ProtobufFieldName:  "room_id",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "start",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "datetime",
			DatabaseTypePretty: "datetime",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "datetime",
			ColumnLength:       -1,
			GoFieldName:        "Start",
			GoFieldType:        "time.Time",
			JSONFieldName:      "start",
			ProtobufFieldName:  "start",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "end",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "datetime",
			DatabaseTypePretty: "datetime",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "datetime",
			ColumnLength:       -1,
			GoFieldName:        "End",
			GoFieldType:        "time.Time",
			JSONFieldName:      "end",
			ProtobufFieldName:  "end",
			ProtobufType:       "google.protobuf.Timestamp",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "title",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "varchar",
			DatabaseTypePretty: "varchar(64)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "varchar",
			ColumnLength:       64,
			GoFieldName:        "Title",
			GoFieldType:        "string",
			JSONFieldName:      "title",
			ProtobufFieldName:  "title",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "event_id",
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
			GoFieldName:        "EventID",
			GoFieldType:        "int32",
			JSONFieldName:      "event_id",
			ProtobufFieldName:  "event_id",
			ProtobufType:       "int32",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "course_code",
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
			GoFieldName:        "CourseCode",
			GoFieldType:        "null.String",
			JSONFieldName:      "course_code",
			ProtobufFieldName:  "course_code",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},
	},
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderSchedules) TableName() string {
	return "roomfinder_schedules"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *RoomfinderSchedules) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (r *RoomfinderSchedules) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (r *RoomfinderSchedules) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (r *RoomfinderSchedules) TableInfo() *TableInfo {
	return roomfinder_schedulesTableInfo
}
