package migration

import (
	"database/sql"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// migrate20210709193000
// adds a "url" column to the database containing the url the file was downloaded from.
// adds a "finished" column to the database that indicates, that a files download is finished.
// adds a "fileDownload" cron job that runs every 5 minutes.
func (m TumDBMigrator) migrate20210709193000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20210709193000",
		Migrate: func(tx *gorm.DB) error {
			type Files struct {
				URL        sql.NullString `gorm:"column:url;default:null;" json:"url"`                         // URL of the file source (if any)
				Downloaded sql.NullBool   `gorm:"column:downloaded;type:boolean;default:1;" json:"downloaded"` // true when file is ready to be served, false when still being downloaded
			}
			if err := tx.AutoMigrate(
				&Files{},
				&model.Crontab{},
			); err != nil {
				return err
			}
			return tx.Create(&model.Crontab{
				Interval: 300,
				Type:     null.String{NullString: sql.NullString{String: "fileDownload", Valid: true}},
			}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropColumn("files", "url"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn("files", "finished"); err != nil {
				return err
			}
			return tx.Delete(&model.Crontab{}, "type = ? AND interval = ?", "fileDownload", 300).Error
		},
	}
}
