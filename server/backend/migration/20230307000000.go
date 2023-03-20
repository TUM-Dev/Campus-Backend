package migration

import (
	"database/sql"
	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

func (m TumDBMigrator) migrate20230307000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230307000000",
		Migrate: func(tx *gorm.DB) error {

			if err := tx.AutoMigrate(
				&model.CrawlerLecture{},
			); err != nil {
				return err
			}

			err := SafeEnumMigrate(tx, model.Crontab{}, "type", "lectureCrawler")
			if err != nil {
				return err
			}

			return tx.Create(&model.Crontab{
				Type: null.String{
					NullString: sql.NullString{
						String: cron.LectureCrawler,
						Valid:  true,
					},
				},
				Interval: 2_628_000,
			}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&model.CrawlerLecture{}); err != nil {
				return err
			}

			err := tx.Delete(&model.Crontab{}, "type = ? AND interval = ?", cron.LectureCrawler, 2_628_000).Error
			if err != nil {
				return err
			}

			return SafeEnumRollback(tx, model.Crontab{}, "type", "lectureCrawler")
		},
	}
}
