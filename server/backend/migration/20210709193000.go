package migration

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

type File struct {
	File       int64       `gorm:"primary_key;autoIncrement;column:file;type:int;" json:"file"`
	Name       string      `gorm:"column:name;type:text;size:16777215;" json:"name"`
	Path       string      `gorm:"column:path;type:text;size:16777215;" json:"path"`
	Downloads  int32       `gorm:"column:downloads;type:int;default:0;" json:"downloads"`
	URL        null.String `gorm:"column:url;default:null;" json:"url"`                         // URL of the files source (if any)
	Downloaded null.Bool   `gorm:"column:downloaded;type:boolean;default:1;" json:"downloaded"` // true when file is ready to be served, false when still being downloaded
}

// migrate20210709193000
// adds a "url" column to the database containing the url the file was downloaded from.
// adds a "finished" column to the database that indicates, that a files download is finished.
// adds a "fileDownload" cron job that runs every 5 minutes.
func migrate20210709193000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20210709193000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(
				&File{},
			); err != nil {
				return err
			}
			type Files struct {
				URL        null.String `gorm:"column:url;default:null;" json:"url"`                         // URL of the file source (if any)
				Downloaded null.Bool   `gorm:"column:downloaded;type:boolean;default:1;" json:"downloaded"` // true when file is ready to be served, false when still being downloaded
			}
			if err := tx.AutoMigrate(
				&Files{},
				&model.Crontab{},
			); err != nil {
				return err
			}
			return tx.Create(&model.Crontab{
				Interval: 300,
				Type:     null.StringFrom("fileDownload"),
			}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropColumn("files", "url"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn("files", "finished"); err != nil {
				return err
			}
			return tx.Delete(&model.Crontab{}, "type = 'fileDownload'").Error
		},
	}
}
