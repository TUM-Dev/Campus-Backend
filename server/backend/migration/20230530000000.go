package migration

import (
	"database/sql"
	_ "embed"
	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

func (m TumDBMigrator) migrate20230530000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230530000000",
		Migrate: func(tx *gorm.DB) error {

			if err := tx.AutoMigrate(
				&model.ExamResultPublished{},
				&model.NewExamResultsSubscriber{},
			); err != nil {
				return err
			}

			err := SafeEnumMigrate(tx, model.Crontab{}, "type", cron.NewExamResultsHook)
			if err != nil {
				return err
			}

			return tx.Create(&model.Crontab{
				Interval: 60, // Every 5 minutes
				Type:     null.StringFrom(cron.NewExamResultsHook),
			}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&model.ExamResultPublished{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&model.NewExamResultsSubscriber{}); err != nil {
				return err
			}

			err := SafeEnumRollback(tx, model.Crontab{}, "type", cron.NewExamResultsHook)
			if err != nil {
				return err
			}

			return tx.Delete(&model.Crontab{}, "type = ?", cron.NewExamResultsHook).Error
		},
	}
}
