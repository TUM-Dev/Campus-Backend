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
			// news_alert does have a FK on the wrong field set
			if err := tx.Exec("alter table news_alert drop foreign key if exists news_alert").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table news_alert drop foreign key if exists news_alert_files_file_fk").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table news_alert add constraint news_alert_files_file_fk foreign key (file) references files (file) on delete cascade").Error; err != nil {
				return err
			}
			// some tables PK don't have AUTO_INCREMENT set, despite that they should
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
			if err := tx.Exec("alter table crontab drop key if exists cron").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table crontab add constraint primary key (cron)").Error; err != nil {
				return err
			}
			if err := migrateField(tx, "crontab", "cron", "BIGINT NOT NULL AUTO_INCREMENT"); err != nil {
				return err
			}
			// roomfinder_schedules does not have a PK set
			if err := tx.Exec("alter table roomfinder_schedules add constraint primary key (room_id, start, end)").Error; err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// news_alert does have a FK on the wrong field set
			if err := tx.Exec("alter table news_alert drop foreign key if exists news_alert_files_file_fk").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table news_alert add constraint news_alert foreign key (news_alert) references files (file) on delete cascade").Error; err != nil {
				return err
			}
			// some tables PK don't have AUTO_INCREMENT set, despite that they should
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
			if err := tx.Exec("alter table crontab drop primary key").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table crontab add constraint cron unique (cron)").Error; err != nil {
				return err
			}
			if err := migrateField(tx, "crontab", "cron", "BIGINT NOT NULL AUTO_INCREMENT"); err != nil {
				return err
			}
			// roomfinder_schedules does not have a PK set
			if err := tx.Exec("alter table roomfinder_schedules drop primary key").Error; err != nil {
				return err
			}
			return nil
		},
	}
}
