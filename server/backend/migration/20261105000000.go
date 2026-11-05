package migration

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20261105000000
// Removes the never-implemented alarm cronjob
func migrate20261105000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20261105000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Delete(&model.Crontab{}, "type = 'alarm'").Error; err != nil {
				return err
			}
			return SafeEnumRemove(tx, &model.Crontab{}, "type", "alarm")
		},
		Rollback: func(tx *gorm.DB) error {
			return SafeEnumAdd(tx, &model.Crontab{}, "type", "alarm")
		},
	}
}
