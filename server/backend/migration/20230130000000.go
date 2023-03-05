package migration

import (
	_ "embed"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func (m TumDBMigrator) migrate20230130000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230130000000",
		Migrate: func(tx *gorm.DB) error {

			if err := tx.AutoMigrate(
				&model.IOSLecture{},
				&model.IOSDeviceLecture{},
				&model.IOSDeviceRequestLog{},
			); err != nil {
				return err
			}

			err := SafeEnumMigrate(tx, model.IOSDeviceRequestLog{}, "request_type", "LECTURE_UPDATE_REQUEST")
			if err != nil {
				return err
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&model.IOSLecture{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&model.IOSDeviceLecture{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&model.IOSDeviceRequestLog{}); err != nil {
				return err
			}

			err := SafeEnumRollback(tx, model.IOSDeviceRequestLog{}, "request_type", "LECTURE_UPDATE_REQUEST")
			if err != nil {
				return err
			}

			return nil
		},
	}
}
