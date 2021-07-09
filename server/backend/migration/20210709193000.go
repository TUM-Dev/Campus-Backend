package migration

import (
	"database/sql"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

//migrate20210709193000
//adds a "url" column to the database containing the url the file was downloaded from.
//adds a "finished" column to the database that indicates, that a files download is finished.
func (m TumDBMigrator) migrate20210709193000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20210709193000",
		Migrate: func(tx *gorm.DB) error {
			type Files struct {
				URL        sql.NullString `gorm:"column:url;default:null;" json:"url"`                       // URL of the file source (if any)
				Downloaded bool           `gorm:"column:finished;type:boolean;default:1;" json:"downloaded"` // true when file is ready to be served, false when still being downloaded
			}
			return tx.AutoMigrate(&Files{})
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropColumn("files", "url"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn("files", "finished"); err != nil {
				return err
			}
			return nil
		},
	}
}
