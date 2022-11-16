package migration

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

//migrate20221115000000

func (m TumDBMigrator) migrate20221116140300() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20221116140300",
		Migrate: func(tx *gorm.DB) error {

			if err := tx.AutoMigrate(
				&model.IOSDevice{},
				&model.IOSDeviceUsageLog{},
			); err != nil {
				return err
			}
			return nil
		},

		Rollback: func(tx *gorm.DB) error {
			return tx.Delete(&model.IOSDevice{}).Error
		},
	}
}
