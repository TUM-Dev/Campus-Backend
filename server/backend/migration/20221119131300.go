package migration

import (
	_ "embed"
	"encoding/json"

	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//migrate20221115000000

//go:embed static_data/iosInitialSchedulingPriorities.json
var iosInitialPrioritiesFile []byte

func (m TumDBMigrator) migrate20221119131300() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20221119131300",
		Migrate: func(tx *gorm.DB) error {

			if err := tx.AutoMigrate(
				&model.IOSDevice{},
				&model.Crontab{},
				&model.IOSSchedulingPriority{},
				&model.IOSScheduledUpdateLog{},
				&model.IOSDeviceRequestLog{},
				&model.IOSEncryptedGrade{},
				&model.IOSDevicesActivityReset{},
			); err != nil {
				return err
			}

			if err := SafeEnumMigrate(tx, &model.Crontab{}, "type", "iosNotifications", "iosActivityReset"); err != nil {
				return err
			}

			var priorities []model.IOSSchedulingPriority

			if err := json.Unmarshal(iosInitialPrioritiesFile, &priorities); err != nil {
				log.WithError(err).Error("could not unmarshal json")
				return err
			}

			if err := tx.Create(&priorities).Error; err != nil {
				log.WithError(err).Error("could not save priority's")
				return err
			}

			err := tx.Create(&model.Crontab{
				Interval: 60,
				Type:     null.StringFrom(cron.IOSNotifications),
			}).Error

			if err != nil {
				log.WithError(err).Error("could not create crontab")
				return err
			}

			return tx.Create(&model.Crontab{
				Type:     null.StringFrom(cron.IOSActivityReset),
				Interval: 86400,
			}).Error
		},

		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&model.IOSDevice{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&model.IOSSchedulingPriority{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&model.IOSScheduledUpdateLog{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&model.IOSDeviceRequestLog{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&model.IOSEncryptedGrade{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&model.IOSDevicesActivityReset{}); err != nil {
				return err
			}

			err := tx.Delete(&model.Crontab{}, "type = ? AND interval = ?", cron.IOSNotifications, 60).Error
			if err != nil {
				return err
			}

			err = tx.Delete(&model.Crontab{}, "type = ? AND interval = ?", cron.IOSActivityReset, 86400).Error

			if err != nil {
				return err
			}

			return SafeEnumRollback(tx, &model.Crontab{}, "type", "iosNotifications", "iosActivityReset")
		},
	}
}
