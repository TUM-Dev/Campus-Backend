package migration

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// migrate20221210000000
// adds a "canteenHeadCount" cron job that runs every 5 minutes.
func migrate20221210000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20221210000000",
		Migrate: func(tx *gorm.DB) error {

			err := tx.AutoMigrate(
				&model.CanteenHeadCount{},
			)
			if err != nil {
				return err
			}

			// allow "canteenHeadCount" in the enum
			if err := SafeEnumAdd(tx, model.Crontab{}, "type", "canteenHeadCount"); err != nil {
				return err
			}

			return tx.Create(&model.Crontab{
				Interval: 60 * 5, // Every 5 minutes
				Type:     null.StringFrom("canteenHeadCount"),
			}).Error
		},

		Rollback: func(tx *gorm.DB) error {
			err := tx.Delete(&model.Crontab{}, "type = 'canteenHeadCount'").Error
			if err != nil {
				return err
			}
			// Remove the 'canteenHeadCount' from the enum
			return SafeEnumRemove(tx, model.Crontab{}, "type", "canteenHeadCount")
		},
	}
}
