package migration

import (
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func tablesWithWrongCOLLATE() []string {
	return []string{"crontab", "devices", "dish", "files", "kino", "news", "newsSource", "notifications", "notification_types"}
}

// migrate20240317000000
// unified a variety of factors to not be different for no reason
func migrate20240317000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240317000000",
		Migrate: func(tx *gorm.DB) error {
			for _, t := range tablesWithWrongCOLLATE() {
				if err := tx.Exec(fmt.Sprintf("ALTER TABLE `%s` COLLATE utf8mb4_unicode_ci", t)).Error; err != nil {
					return err
				}
			}
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
