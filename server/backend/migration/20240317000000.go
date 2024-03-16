package migration

import (
	"fmt"
	"os"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func tablesWithWrongCOLLATE() []string {
	return []string{"crontab", "devices", "dish", "files", "kino", "news", "newsSource", "notification", "notification_type", "notification_confirmation", "feedback", "update_note", "news_alert"}
}

type columnsWithWrongCollationOrCharset struct {
	tableName        string
	columnName       string
	columnType       string
	characterSetName string
	collationName    string
}

// feedbackColumnsWithWrongCOLLATE lists all columns that need changing
// can be gotten with
// SELECT TABLE_NAME, COLUMN_NAME, COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME
// from information_schema.columns
// WHERE TABLE_SCHEMA = 'campus_db'
//
//	and (CHARACTER_SET_NAME != 'utf8mb4' or COLLATION_NAME != 'utf8mb4_unicode_ci')
func feedbackColumnsWithWrongCOLLATE() []columnsWithWrongCollationOrCharset {
	return []columnsWithWrongCollationOrCharset{
		{"barrierFree_moreInfo", "title", "varchar(32)", "utf8mb3", "utf8mb3_general_ci"},
		{"barrierFree_moreInfo", "category", "varchar(11)", "utf8mb3", "utf8mb3_general_ci"},
		{"barrierFree_moreInfo", "url", "varchar(128)", "utf8mb3", "utf8mb3_general_ci"},
		{"barrierFree_persons", "name", "varchar(40)", "utf8mb3", "utf8mb3_general_ci"},
		{"barrierFree_persons", "telephone", "varchar(32)", "utf8mb3", "utf8mb3_general_ci"},
		{"barrierFree_persons", "email", "varchar(32)", "utf8mb3", "utf8mb3_general_ci"},
		{"barrierFree_persons", "faculty", "varchar(32)", "utf8mb3", "utf8mb3_general_ci"},
		{"barrierFree_persons", "office", "varchar(16)", "utf8mb3", "utf8mb3_general_ci"},
		{"barrierFree_persons", "officeHour", "varchar(16)", "utf8mb3", "utf8mb3_general_ci"},
		{"barrierFree_persons", "tumID", "varchar(24)", "utf8mb3", "utf8mb3_general_ci"},
		{"event", "title", "varchar(100)", "utf8mb3", "utf8mb3_general_ci"},
		{"event", "description", "text", "utf8mb3", "utf8mb3_general_ci"},
		{"event", "locality", "varchar(200)", "utf8mb3", "utf8mb3_general_ci"},
		{"event", "link", "varchar(200)", "utf8mb3", "utf8mb3_general_ci"},
		{"faculty", "name", "varchar(150)", "utf8mb4", "utf8mb4_general_ci"},
		{"feedback", "email_id", "text", "utf8mb3", "utf8mb3_general_ci"},
		{"feedback", "receiver", "text", "utf8mb3", "utf8mb3_general_ci"},
		{"feedback", "reply_to", "text", "utf8mb3", "utf8mb3_general_ci"},
		{"feedback", "feedback", "text", "utf8mb3", "utf8mb3_general_ci"},
		{"feedback", "os_version", "text", "utf8mb4", "utf8mb4_general_ci"},
		{"feedback", "app_version", "text", "utf8mb4", "utf8mb4_general_ci"},
		{"location", "name", "text", "utf8mb3", "utf8mb3_general_ci"},
		{"news_alert", "name", "varchar(100)", "utf8mb4", "utf8mb4_general_ci"},
		{"news_alert", "link", "text", "utf8mb4", "utf8mb4_general_ci"},
		{"notification", "title", "text", "utf8mb3", "utf8mb3_general_ci"},
		{"notification", "description", "text", "utf8mb3", "utf8mb3_general_ci"},
		{"notification", "signature", "text", "utf8mb3", "utf8mb3_general_ci"},
		{"notification_type", "name", "text", "utf8mb3", "utf8mb3_general_ci"},
		{"notification_type", "confirmation", "enum('true','false')", "utf8mb3", "utf8mb3_general_ci"},
		{"question", "text", "text", "utf8mb4", "utf8mb4_general_ci"},
		{"questionAnswers", "text", "text", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_building2area", "building_nr", "varchar(8)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_building2area", "campus", "char(1)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_building2area", "name", "varchar(32)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_buildings", "building_nr", "varchar(8)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_buildings", "utm_zone", "varchar(4)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_buildings", "utm_easting", "varchar(32)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_buildings", "utm_northing", "varchar(32)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_buildings2gps", "id", "varchar(8)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_buildings2gps", "latitude", "varchar(30)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_buildings2gps", "longitude", "varchar(30)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_buildings2maps", "building_nr", "varchar(8)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_maps", "description", "varchar(64)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_rooms", "room_code", "varchar(32)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_rooms", "building_nr", "varchar(8)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_rooms", "arch_id", "varchar(16)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_rooms", "info", "varchar(64)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_rooms", "address", "varchar(128)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_rooms", "purpose", "varchar(64)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_rooms", "utm_zone", "varchar(4)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_rooms", "utm_easting", "varchar(32)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_rooms", "utm_northing", "varchar(32)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_schedules", "title", "varchar(64)", "utf8mb4", "utf8mb4_general_ci"},
		{"roomfinder_schedules", "course_code", "varchar(32)", "utf8mb4", "utf8mb4_general_ci"},
		{"sessions", "session", "varchar(255)", "utf8mb3", "utf8mb3_general_ci"},
		{"ticket_admin", "key", "text", "utf8mb3", "utf8mb3_general_ci"},
		{"ticket_admin", "comment", "text", "utf8mb3", "utf8mb3_general_ci"},
		{"ticket_group", "description", "text", "utf8mb3", "utf8mb3_general_ci"},
		{"ticket_history", "code", "char(128)", "utf8mb3", "utf8mb3_general_ci"},
		{"ticket_payment", "name", "varchar(50)", "utf8mb3", "utf8mb3_general_ci"},
		{"ticket_payment", "config", "text", "utf8mb3", "utf8mb3_general_ci"},
		{"ticket_type", "description", "varchar(100)", "utf8mb3", "utf8mb3_general_ci"},
		{"update_note", "version_name", "text", "utf8mb4", "utf8mb4_general_ci"},
		{"update_note", "message", "text", "utf8mb4", "utf8mb4_general_ci"},
	}
}

// migrate20240317000000
// unified all of our tables, the database and the fields to use `utf8mb4_unicode_ci` instead of the legacy `utf8mb3_general_ci` or `latin1`
func migrate20240317000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240317000000",
		Migrate: func(tx *gorm.DB) error {
			// first migrate the db
			if err := tx.Exec(fmt.Sprintf("ALTER DATABASE `%s` CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci", os.Getenv("DB_NAME"))).Error; err != nil {
				return err
			}
			// then set the tables
			for _, t := range tablesWithWrongCOLLATE() {
				if err := tx.Exec(fmt.Sprintf("ALTER TABLE `%s` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", t)).Error; err != nil {
					return err
				}
			}
			// then convert single columns in each table
			for _, f := range feedbackColumnsWithWrongCOLLATE() {
				if err := tx.Exec(fmt.Sprintf("alter table %s modify %s %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", f.tableName, f.columnName, f.columnType)).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// first migrate the db
			if err := tx.Exec(fmt.Sprintf("ALTER DATABASE `%s` CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci", os.Getenv("DB_NAME"))).Error; err != nil {
				return err
			}
			// revert changes to tables
			for _, t := range tablesWithWrongCOLLATE() {
				if err := tx.Exec(fmt.Sprintf("ALTER TABLE `%s` COLLATE utf8mb4_general_ci", t)).Error; err != nil {
					return err
				}
			}
			// revert changes to fields
			for _, f := range feedbackColumnsWithWrongCOLLATE() {
				if err := tx.Exec(fmt.Sprintf("alter table %s modify %s %s CHARACTER SET %s COLLATE %s", f.tableName, f.columnName, f.columnType, f.characterSetName, f.collationName)).Error; err != nil {
					return err
				}
			}
			return nil
		},
	}
}
