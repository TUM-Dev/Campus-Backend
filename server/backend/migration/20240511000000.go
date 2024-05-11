package migration

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

type device2statsRooms struct {
	RoomFinderActivity        int `gorm:"column:RoomFinderActivity;default 0;not null"`
	RoomFinderDetailsActivity int `gorm:"column:RoomFinderDetailsActivity;default 0;not null"`
}

// TableName sets the insert table name for this struct type
func (n *device2statsRooms) TableName() string {
	return "device2stats"
}

// migrate20240511000000
// - Removes all traces of the room-finder from the database
func migrate20240511000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240511000000",
		Migrate: func(tx *gorm.DB) error {
			// remove the cronjob
			if err := SafeEnumRemove(tx, &model.Crontab{}, "type", "roomfinder"); err != nil {
				return err
			}
			if err := tx.Delete(&model.Crontab{}, "type = 'roomfinder'").Error; err != nil {
				return err
			}
			// Remove tracking from device2stats
			if err := tx.Migrator().DropColumn(&device2statsRooms{}, "RoomFinderActivity"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&device2statsRooms{}, "RoomFinderDetailsActivity"); err != nil {
				return err
			}
			// remove redundant tables
			tablesToBeDeleted := []string{
				"roomfinder_buildings",
				"roomfinder_buildings2gps",
				"roomfinder_buildings2maps",
				"roomfinder_maps",
				"roomfinder_rooms",
				"roomfinder_rooms2maps",
				"roomfinder_schedules",
				"roomfinder_building2area",
			}
			for _, table := range tablesToBeDeleted {
				if err := tx.Migrator().DropTable(table); err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// re-add the cronjob
			if err := SafeEnumAdd(tx, &model.Crontab{}, "type", "roomfinder"); err != nil {
				return err
			}
			for i := 0; i <= 10; i++ {
				if err := tx.Create(&model.Crontab{
					Interval: 10 * 60 * 60 * 24, // Every 10 days
					Type:     null.StringFrom("roomfinder"),
					ID:       null.IntFrom(int64(i)),
				}).Error; err != nil {
					return err
				}
			}

			// Add tracking to device2stats
			if err := tx.Migrator().AutoMigrate(&device2statsRooms{}); err != nil {
				return err
			}
			// re-add redundant tables
			tableDefinition := []string{
				"create table roomfinder_buildings ( building_nr varchar(8) collate utf8mb4_unicode_ci not null primary key, utm_zone varchar(4) collate utf8mb4_unicode_ci null, utm_easting varchar(32) collate utf8mb4_unicode_ci null, utm_northing varchar(32) collate utf8mb4_unicode_ci null, default_map_id int null ) charset = utf8mb4;",
				"create table roomfinder_buildings2gps ( id varchar(8) collate utf8mb4_unicode_ci not null primary key, latitude varchar(30) collate utf8mb4_unicode_ci null, longitude varchar(30) collate utf8mb4_unicode_ci null ) charset = utf8mb4;",
				"create table roomfinder_buildings2maps ( building_nr varchar(8) collate utf8mb4_unicode_ci not null, map_id bigint not null, primary key (building_nr, map_id) ) charset = utf8mb4;",
				"create table roomfinder_maps ( map_id bigint auto_increment primary key, description varchar(64) collate utf8mb4_unicode_ci null, scale int not null, width int not null, height int not null ) charset = utf8mb4;",
				"create table roomfinder_rooms ( room_id bigint auto_increment primary key, room_code varchar(32) collate utf8mb4_unicode_ci null, building_nr varchar(8) collate utf8mb4_unicode_ci null, arch_id varchar(16) collate utf8mb4_unicode_ci null, info varchar(64) collate utf8mb4_unicode_ci null, address varchar(128) collate utf8mb4_unicode_ci null, purpose_id bigint null, purpose varchar(64) collate utf8mb4_unicode_ci null, seats bigint null, utm_zone varchar(4) collate utf8mb4_unicode_ci null, utm_easting varchar(32) collate utf8mb4_unicode_ci null, utm_northing varchar(32) collate utf8mb4_unicode_ci null, unit_id bigint null, default_map_id bigint null ) charset = utf8mb4;",
				"create table roomfinder_rooms2maps ( room_id bigint not null, map_id bigint not null, primary key (room_id, map_id) ) charset = utf8mb4;",
				"create table roomfinder_schedules ( room_id bigint auto_increment, start datetime not null, end datetime not null, title varchar(64) collate utf8mb4_unicode_ci null, event_id int not null, course_code varchar(32) collate utf8mb4_unicode_ci null, primary key (room_id, start, end), constraint `unique` unique (room_id, start, end) ) charset = utf8mb4;",
				"create table roomfinder_building2area(area_id int not null, building_nr varchar(8) collate utf8mb4_unicode_ci not null primary key, campus char collate utf8mb4_unicode_ci null, name varchar(32) collate utf8mb4_unicode_ci null) charset = utf8mb4;",
			}
			for _, defintion := range tableDefinition {
				if err := tx.Exec(defintion).Error; err != nil {
					return err
				}
			}
			return nil
		},
	}
}
