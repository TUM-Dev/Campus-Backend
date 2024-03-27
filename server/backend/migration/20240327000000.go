package migration

import (
	"slices"

	"github.com/go-gormigrate/gormigrate/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// migrate20240327000000
// made sure auto_increment is re-added to all fields where this was accidentally removed in migrate20240316000000
func migrate20240327000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240327000000",
		Migrate: func(tx *gorm.DB) error {
			for _, t := range tablesWithWrongId() {
				tablesWithCorrectIds := []string{"roomfinder_buildings2maps", "roomfinder_rooms2maps", "question"}
				if slices.Contains(tablesWithCorrectIds, t.table) {
					continue
				}
				log.WithField("table", t.table).Info("migrated PK-field")
				if err := migrateField(tx, t.table, t.field, "BIGINT NOT NULL AUTO_INCREMENT"); err != nil {
					return err
				}
			}
			// crontab should be a PK instead of a unique key
			if err := migrateField(tx, "crontab", "cron", "BIGINT NOT NULL"); err != nil {
				return err
			}
			if err := tx.Exec("alter table crontab drop key cron").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table crontab add constraint crontab_pk primary key (cron)").Error; err != nil {
				return err
			}
			if err := migrateField(tx, "crontab", "cron", "BIGINT NOT NULL AUTO_INCREMENT"); err != nil {
				return err
			}
			// roomfinder_schedules does not have a PK set
			if err := tx.Exec("alter table roomfinder_schedules add constraint roomfinder_schedules_pk primary key (room_id)").Error; err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			for _, t := range tablesWithWrongId() {
				tablesWithCorrectIds := []string{"roomfinder_buildings2maps", "roomfinder_rooms2maps", "question"}
				if slices.Contains(tablesWithCorrectIds, t.table) {
					continue
				}
				if err := migrateField(tx, t.table, t.field, "BIGINT NOT NULL"); err != nil {
					return err
				}
				log.WithField("table", t.table).Info("migrated PK-field")
			}
			// crontab should be a PK instead of a unique key
			if err := migrateField(tx, "crontab", "cron", "BIGINT NOT NULL"); err != nil {
				return err
			}
			if err := tx.Exec("alter table crontab drop key crontab_pk").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table crontab add constraint cron unique (cron)").Error; err != nil {
				return err
			}
			if err := migrateField(tx, "crontab", "cron", "BIGINT NOT NULL AUTO_INCREMENT"); err != nil {
				return err
			}
			// roomfinder_schedules does not have a PK set
			if err := tx.Exec("alter table roomfinder_schedules drop key roomfinder_schedules_pk").Error; err != nil {
				return err
			}
			return nil
		},
	}
}
