package migration

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// NewsSource struct is a row record of the newsSource table in the tca database
type NewsSource struct {
	Source int64       `gorm:"primary_key;AUTO_INCREMENT;column:source;type:int;"`
	Title  string      `gorm:"column:title;type:text;size:16777215;"`
	URL    null.String `gorm:"column:url;type:text;size:16777215;"`
	FileID int64       `gorm:"column:icon;not null;type:int;"`
	File   model.File  `gorm:"foreignKey:FileID;references:file;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Hook   null.String `gorm:"column:hook;type:char;size:12;"`
}

// migrate20230904100000
// migrates the crontab from kino to movie crontab
func migrate20230904100000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230904100000",
		Migrate: func(tx *gorm.DB) error {
			// modify the crontab
			if err := tx.Delete(&model.Crontab{}, "type = 'kino'").Error; err != nil {
				return err
			}
			if err := SafeEnumRemove(tx, &model.Crontab{}, "type", "kino"); err != nil {
				return err
			}
			if err := SafeEnumAdd(tx, &model.Crontab{}, "type", "movie"); err != nil {
				return err
			}
			if err := tx.AutoMigrate(&NewsSource{}); err != nil {
				return err
			}
			// tu film news source is now inlined
			if err := tx.Delete(&NewsSource{Source: 2}).Error; err != nil {
				return err
			}
			return tx.Create(&model.Crontab{
				Interval: 60 * 60 * 24, // daily
				Type:     null.StringFrom("movie"),
			}).Error
		},

		Rollback: func(tx *gorm.DB) error {
			// modify the crontab
			if err := tx.Delete(&model.Crontab{}, "type = 'movie'").Error; err != nil {
				return err
			}
			if err := SafeEnumRemove(tx, &model.Crontab{}, "type", "movie"); err != nil {
				return err
			}
			if err := SafeEnumAdd(tx, &model.Crontab{}, "type", "kino"); err != nil {
				return err
			}
			if err := tx.Create(&NewsSource{
				Source: 2,
				Title:  "TU Film",
				URL:    null.StringFrom("http://www.tu-film.de/programm/index/upcoming.rss"),
				FileID: 2,
				Hook:   null.String{},
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
