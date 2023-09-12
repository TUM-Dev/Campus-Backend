package migration

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20230825000000
// Removes the ability to run chat cronjobs
func (m TumDBMigrator) migrate20230825000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230825000000",
		Migrate: func(tx *gorm.DB) error {
			// deactivete the crontab (Rollback deletes this from the enum)
			// given that previously, not cronjobs for this type existed there is no need to remove offending entries first
			return SafeEnumRollback(tx, &model.Crontab{}, "type", "chat")
		},
		Rollback: func(tx *gorm.DB) error {
			// activete the crontab (Migrate adds this from to the enum)
			// given that previously, not cronjobs for this type existed there is no need to add entries first
			return SafeEnumMigrate(tx, &model.Crontab{}, "type", "chat")
		},
	}
}
