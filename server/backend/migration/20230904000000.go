package migration

import (
	"database/sql"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// migrate20230904000000
// Removes ticketsales from the db-enums
func (m TumDBMigrator) migrate20230904000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230904000000",
		Migrate: func(tx *gorm.DB) error {
			// remove "canteenHeadCount" in the enum
			if err := tx.Delete(&model.Crontab{}, "type = 'ticketsales'").Error; err != nil {
				return err
			}
			if err := SafeEnumRollback(tx, model.Crontab{}, "type", "ticketsales"); err != nil {
				return err
			}
			return nil
		},

		Rollback: func(tx *gorm.DB) error {
			if err := SafeEnumMigrate(tx, model.Crontab{}, "type", "ticketsales"); err != nil {
				return err
			}
			return tx.Create(&model.Crontab{
				Interval: 60 * 10, // Every 10 minutes
				Type:     null.String{NullString: sql.NullString{String: "ticketsales", Valid: true}},
			}).Error
		},
	}
}
