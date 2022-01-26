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

// RoomfinderBuildings2gps struct is a row record of the roomfinder_buildings2gps table in the tca database
type RoomfinderBuildings2gps struct {
	//[ 0] id                                             varchar(8)           null: false  primary: true   isArray: false  auto: false  col: varchar         len: 8       default: []
	ID string `gorm:"primary_key;column:id;type:varchar;size:8;" json:"id"`
	//[ 1] latitude                                       varchar(30)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 30      default: []
	Latitude null.String `gorm:"column:latitude;type:varchar;size:30;" json:"latitude"`
	//[ 2] longitude                                      varchar(30)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 30      default: []
	Longitude null.String `gorm:"column:longitude;type:varchar;size:30;" json:"longitude"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderBuildings2gps) TableName() string {
	return "roomfinder_buildings2gps"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (r *RoomfinderBuildings2gps) BeforeSave() error {
	return nil
}
