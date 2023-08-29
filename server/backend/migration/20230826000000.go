package migration

import (
	"database/sql"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Feedback struct {
	OsVersion  sql.NullString `gorm:"column:os_version;type:text;null;"`
	AppVersion sql.NullString `gorm:"column:app_version;type:text;null;"`
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
			if err := tx.Migrator().AddColumn(&Feedback{}, "OsVersion"); err != nil {
				return err
			}
			return tx.Migrator().AddColumn(&Feedback{}, "AppVersion")
		},

		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropColumn(&Feedback{}, "OsVersion"); err != nil {
				return err
			}
			return tx.Migrator().DropColumn(&Feedback{}, "AppVersion")
		},
	}
}
