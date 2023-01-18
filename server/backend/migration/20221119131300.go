package migration

import (
	"database/sql"
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
				&model.IOSLog{},
				&model.IOSDevicesActivityReset{},
			); err != nil {
				return err
			}

			err := SafeEnumMigrate(tx, &model.Crontab{}, "type", "iosNotifications", "iosActivityReset")

			if err != nil {
				return err
			}

			var priorities []model.IOSSchedulingPriority

			unmarshalErr := json.Unmarshal(iosInitialPrioritiesFile, &priorities)

			if unmarshalErr != nil {
				log.Info(unmarshalErr.Error())
				return unmarshalErr
			}

			if err := tx.Create(&priorities).Error; err != nil {
				log.Info(err.Error())
				return err
			}

			err = tx.Create(&model.Crontab{
				Interval: 60,
				Type:     null.String{NullString: sql.NullString{String: cron.IOSNotifications, Valid: true}},
			}).Error

			if err != nil {
				log.Error(err.Error())
				return err
			}

			return tx.Create(&model.Crontab{
				Type: null.String{
					NullString: sql.NullString{
						String: cron.IOSActivityReset,
						Valid:  true,
					},
				},
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
			if err := tx.Migrator().DropTable(&model.IOSLog{}); err != nil {
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
