package migration

import (
	"database/sql"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// migrate20221210000000
// adds a "canteenHeadCount" cron job that runs every 5 minutes.
func (m TumDBMigrator) migrate20221210000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20221210000000",
		Migrate: func(tx *gorm.DB) error {

			if err := tx.AutoMigrate(
				&model.CanteenHeadCount{},
			); err != nil {
				return err
			}

			return tx.Create(&model.Crontab{
				Interval: 60 * 5, // Every 5 minutes
				Type:     null.String{NullString: sql.NullString{String: "canteenHeadCount", Valid: true}},
			}).Error
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Delete(&model.Crontab{}, "type = ? AND interval = ?", "canteenHeadCount", 60*5).Error
		},
	}
}
