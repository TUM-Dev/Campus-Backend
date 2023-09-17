package migration

import (
	"database/sql"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// migrate20230825000000
// Removes the ability to run chat cronjobs
func (m TumDBMigrator) migrate20230825000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230825000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Delete(&model.Crontab{}, "type = 'chat'").Error; err != nil {
				return err
			}
			return SafeEnumRollback(tx, &model.Crontab{}, "type", "chat")
		},
		Rollback: func(tx *gorm.DB) error {
			if err := SafeEnumMigrate(tx, &model.Crontab{}, "type", "chat"); err != nil {
				return err
			}
			return tx.Create(&model.Crontab{
				Interval: 60 * 10, // Every 10 minutes
				Type:     null.String{NullString: sql.NullString{String: "chat", Valid: true}},
			}).Error
		},
	}
}
