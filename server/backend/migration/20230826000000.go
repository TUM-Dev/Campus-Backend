package migration

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

type Feedback struct {
	Processed  bool        `gorm:"column:processed;type:boolean;default:false;not null;"`
	OsVersion  null.String `gorm:"column:os_version;type:text;null;"`
	AppVersion null.String `gorm:"column:app_version;type:text;null;"`
}

// TableName sets the insert table name for this struct type
func (n *Feedback) TableName() string {
	return "feedback"
}

// migrate20230826000000
// adds a "feedbackEmail" cron job that runs every 30 minutes.
func (m TumDBMigrator) migrate20230826000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20230826000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Migrator().AddColumn(&Feedback{}, "Processed"); err != nil {
				return err
			}
			if err := tx.Migrator().AddColumn(&Feedback{}, "OsVersion"); err != nil {
				return err
			}
			if err := tx.Migrator().AddColumn(&Feedback{}, "AppVersion"); err != nil {
				return err
			}
			if err := tx.Exec("UPDATE feedback SET processed = true WHERE processed != true;").Error; err != nil {
				return err
			}
			if err := SafeEnumMigrate(tx, &model.Crontab{}, "type", "feedbackEmail"); err != nil {
				return err
			}
			return tx.Create(&model.Crontab{
				Interval: 60 * 30, // Every 30 minutes
				Type:     null.StringFrom("feedbackEmail"),
			}).Error
		},

		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropColumn(&Feedback{}, "Processed"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&Feedback{}, "OsVersion"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&Feedback{}, "AppVersion"); err != nil {
				return err
			}
			if err := tx.Delete(&model.Crontab{Type: null.StringFrom("fileDownload")}).Error; err != nil {
				return err
			}
			return SafeEnumMigrate(tx, &model.Crontab{}, "type", "feedbackEmail")
		},
	}
}
