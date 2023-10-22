package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// RoomfinderRooms struct is a row record of the roomfinder_rooms table in the tca database
type RoomfinderRooms struct {
	RoomID       int64       `gorm:"primary_key;column:room_id;type:int;" json:"room_id"`
	RoomCode     null.String `gorm:"column:room_code;type:varchar(32);" json:"room_code"`
	BuildingNr   null.String `gorm:"column:building_nr;type:varchar(8);" json:"building_nr"`
	ArchID       null.String `gorm:"column:arch_id;type:varchar(16);" json:"arch_id"`
	Info         null.String `gorm:"column:info;type:varchar(64);" json:"info"`
	Address      null.String `gorm:"column:address;type:varchar(128);" json:"address"`
	PurposeID    null.Int    `gorm:"column:purpose_id;type:int;" json:"purpose_id"`
	Purpose      null.String `gorm:"column:purpose;type:varchar(64);" json:"purpose"`
	Seats        null.Int    `gorm:"column:seats;type:int;" json:"seats"`
	UtmZone      null.String `gorm:"column:utm_zone;type:varchar(4);" json:"utm_zone"`
	UtmEasting   null.String `gorm:"column:utm_easting;type:varchar(32);" json:"utm_easting"`
	UtmNorthing  null.String `gorm:"column:utm_northing;type:varchar(32);" json:"utm_northing"`
	UnitID       null.Int    `gorm:"column:unit_id;type:int;" json:"unit_id"`
	DefaultMapID null.Int    `gorm:"column:default_map_id;type:int;" json:"default_map_id"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderRooms) TableName() string {
	return "roomfinder_rooms"
}

// RoomfinderBuilding2area struct is a row record of the roomfinder_building2area table in the tca database
type RoomfinderBuilding2area struct {
	BuildingNr string `gorm:"primary_key;column:building_nr;type:varchar(8);" json:"building_nr"`
	AreaID     int32  `gorm:"column:area_id;type:int;" json:"area_id"`
	Campus     string `gorm:"column:campus;type:char;size:1;" json:"campus"`
	Name       string `gorm:"column:name;type:varchar(32);" json:"name"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderBuilding2area) TableName() string {
	return "roomfinder_building2area"
}

// RoomfinderBuildings struct is a row record of the roomfinder_buildings table in the tca database
type RoomfinderBuildings struct {
	BuildingNr   string      `gorm:"primary_key;column:building_nr;type:varchar(8);" json:"building_nr"`
	UtmZone      null.String `gorm:"column:utm_zone;type:varchar(4);" json:"utm_zone"`
	UtmEasting   null.String `gorm:"column:utm_easting;type:varchar(32);" json:"utm_easting"`
	UtmNorthing  null.String `gorm:"column:utm_northing;type:varchar(32);" json:"utm_northing"`
	DefaultMapID null.Int    `gorm:"column:default_map_id;type:int;" json:"default_map_id"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderBuildings) TableName() string {
	return "roomfinder_buildings"
}

// RoomfinderBuildings2gps struct is a row record of the roomfinder_buildings2gps table in the tca database
type RoomfinderBuildings2gps struct {
	ID        string      `gorm:"primary_key;column:id;type:varchar(8);" json:"id"`
	Latitude  null.String `gorm:"column:latitude;type:varchar(30);" json:"latitude"`
	Longitude null.String `gorm:"column:longitude;type:varchar(30);" json:"longitude"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderBuildings2gps) TableName() string {
	return "roomfinder_buildings2gps"
}

// RoomfinderBuildings2maps struct is a row record of the roomfinder_buildings2maps table in the tca database
type RoomfinderBuildings2maps struct {
	BuildingNr string `gorm:"primary_key;column:building_nr;type:varchar(8);" json:"building_nr"`
	MapID      int64  `gorm:"primary_key;column:map_id;type:int;" json:"map_id"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderBuildings2maps) TableName() string {
	return "roomfinder_buildings2maps"
}

// RoomfinderRooms2maps struct is a row record of the roomfinder_rooms2maps table in the tca database
type RoomfinderRooms2maps struct {
	RoomID int64 `gorm:"primary_key;column:room_id;type:int;" json:"room_id"`
	MapID  int64 `gorm:"primary_key;column:map_id;type:int;" json:"map_id"`
}

// TableName sets the insert table name for this struct type
func (r *RoomfinderRooms2maps) TableName() string {
	return "roomfinder_rooms2maps"
}

// migrate20220126230000
// adds a fulltext index to the roomfinder_rooms table
func (m TumDBMigrator) migrate20220126230000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20220126230000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(
				&RoomfinderRooms{},
				&RoomfinderBuilding2area{},
				&RoomfinderBuildings2gps{},
				&RoomfinderBuildings2maps{},
				&RoomfinderRooms2maps{},
			); err != nil {
				return err
			}

			return tx.Exec("CREATE FULLTEXT INDEX `search_index` ON `roomfinder_rooms` (`info`, `address`, `room_code`)").Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec("DROP INDEX search_index ON roomfinder_rooms").Error
		},
	}
}
