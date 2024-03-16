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
			// todo: then convert single columns in each table
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			for _, t := range tablesWithWrongCOLLATE() {
				if err := tx.Exec(fmt.Sprintf("ALTER TABLE `%s` COLLATE utf8mb4_general_ci", t)).Error; err != nil {
					return err
				}
			}
			return nil
		},
	}
}
