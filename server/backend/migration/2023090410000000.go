package migration

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// migrate2023090410000000
// migrates the crontap from kino to movie crontab
func (m TumDBMigrator) migrate2023090410000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2023090410000000",
		Migrate: func(tx *gorm.DB) error {
			// modify the crontab
			if err := tx.Delete(&model.Crontab{}, "type = ?", "kino").Error; err != nil {
				return err
			}
			if err := SafeEnumRollback(tx, &model.Crontab{}, "type", "kino"); err != nil {
				return err
			}
			if err := SafeEnumMigrate(tx, &model.Crontab{}, "type", "movie"); err != nil {
				return err
			}
			// tu film news source is now inlined
			if err := tx.Delete(&model.NewsSource{Source: 2}).Error; err != nil {
				return err
			}
			return tx.Create(&model.Crontab{
				Interval: 60 * 60 * 24, // daily
				Type:     null.StringFrom("movie"),
			}).Error
		},

		Rollback: func(tx *gorm.DB) error {
			// modify the crontab
			if err := tx.Delete(&model.Crontab{}, "type = ?", "movie").Error; err != nil {
				return err
			}
			if err := SafeEnumRollback(tx, &model.Crontab{}, "type", "movie"); err != nil {
				return err
			}
			if err := SafeEnumMigrate(tx, &model.Crontab{}, "type", "kino"); err != nil {
				return err
			}
			if err := tx.Create(&model.NewsSource{
				Source:  2,
				Title:   "TU Film",
				URL:     null.StringFrom("http://www.tu-film.de/programm/index/upcoming.rss"),
				FilesID: 2,
				Hook:    null.String{},
			}).Error; err != nil {
				return err
			}
			return tx.Create(&model.Crontab{
				Interval: 24 * 60 * 60, // daily
				Type:     null.StringFrom("kino"),
				ID:       null.IntFrom(2),
			}).Error
		},
	}
}
