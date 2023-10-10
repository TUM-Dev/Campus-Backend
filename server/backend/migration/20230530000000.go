package migration

import (
	_ "embed"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

type PublishedExamResult struct {
	Date         time.Time
	ExamID       string `gorm:"primary_key"`
	LectureTitle string
	LectureType  string
	LectureSem   string
	Published    bool
}

type NewExamResultsSubscriber struct {
	CallbackUrl    string `gorm:"primary_key"`
	ApiKey         null.String
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	LastNotifiedAt null.Time
}

func (m TumDBMigrator) migrate20230530000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230530000000",
		Migrate: func(tx *gorm.DB) error {

			if err := tx.AutoMigrate(
				&PublishedExamResult{},
				&NewExamResultsSubscriber{},
			); err != nil {
				return err
			}

			err := SafeEnumAdd(tx, model.Crontab{}, "type", "newExamResultsHook")
			if err != nil {
				return err
			}

			return tx.Create(&model.Crontab{
				Interval: 60 * 10, // Every 10 minutes
				Type:     null.StringFrom("newExamResultsHook"),
			}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&PublishedExamResult{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable(&NewExamResultsSubscriber{}); err != nil {
				return err
			}

			err := SafeEnumRemove(tx, model.Crontab{}, "type", "newExamResultsHook")
			if err != nil {
				return err
			}

			return tx.Delete(&model.Crontab{}, "type = 'newExamResultsHook'").Error
		},
	}
}
