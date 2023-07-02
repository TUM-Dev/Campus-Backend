package migration

import (
	"database/sql"
	_ "embed"
	"github.com/TUM-Dev/Campus-Backend/server/env"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (m TumDBMigrator) migrate20230618000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230618000000",
		Migrate: func(tx *gorm.DB) error {

			if err := tx.AutoMigrate(
				&model.DbExam{},
				&model.DeviceExam{},
				&model.EncryptedGrade{},
			); err != nil {
				return err
			}

			err := tx.Migrator().DropColumn(&model.IOSDevice{}, "activity_today")
			if err != nil {
				log.WithError(err).Info("Could not drop column activity_today")
				return err
			}
			err = tx.Migrator().DropColumn(&model.IOSDevice{}, "activity_this_week")
			if err != nil {
				log.WithError(err).Info("Could not drop column activity_this_week")
				return err
			}
			err = tx.Migrator().DropColumn(&model.IOSDevice{}, "activity_this_month")
			if err != nil {
				log.WithError(err).Info("Could not drop column activity_this_month")
				return err
			}
			err = tx.Migrator().DropColumn(&model.IOSDevice{}, "activity_this_year")
			if err != nil {
				log.WithError(err).Info("Could not drop column activity_this_year")
				return err
			}

			err = tx.Create(&model.NewExamResultsSubscriber{
				CallbackUrl: env.ApiUrl(),
				ApiKey: sql.NullString{
					String: env.ApiKey(),
					Valid:  true,
				},
			}).Error
			if err != nil {
				log.WithError(err).Info("Could not create new exam results subscriber")
				return err
			}

			err = SafeEnumRollback(tx, &model.Crontab{}, "type", "iosNotifications", "iosActivityReset")

			return err
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&model.DbExam{}); err != nil {
				return err
			}

			if err := tx.Migrator().DropTable(&model.DeviceExam{}); err != nil {
				return err
			}

			err := tx.Migrator().AddColumn(&model.IOSDevice{}, "activity_today")
			if err != nil {
				log.WithError(err).Info("Could not add column activity_today")
				return err
			}
			err = tx.Migrator().AddColumn(&model.IOSDevice{}, "activity_this_week")
			if err != nil {
				log.WithError(err).Info("Could not add column activity_this_week")
				return err
			}
			err = tx.Migrator().AddColumn(&model.IOSDevice{}, "activity_this_month")
			if err != nil {
				log.WithError(err).Info("Could not add column activity_this_month")
				return err
			}
			err = tx.Migrator().AddColumn(&model.IOSDevice{}, "activity_this_year")
			if err != nil {
				log.WithError(err).Info("Could not add column activity_this_year")
				return err
			}

			err = tx.Delete(&model.NewExamResultsSubscriber{}, "callback_url = ?", env.ApiUrl()).Error
			if err != nil {
				log.WithError(err).Info("Could not delete new exam results subscriber")
				return err
			}

			err = SafeEnumMigrate(tx, &model.Crontab{}, "type", "iosNotifications", "iosActivityReset")
			if err != nil {
				log.WithError(err).Info("Could not migrate crontab type enum")
				return err
			}

			if err != nil {
				return err
			}

			return nil
		},
	}
}
