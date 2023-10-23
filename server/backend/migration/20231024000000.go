package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type wrongTableName struct {
	Original string
	New      string
}

var wrongTableNames = []wrongTableName{}

// migrate20231024000000
// - replaces all instances of misleadingly named tables with the correct ones
func migrate20231024000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20231024000000",
		Migrate: func(tx *gorm.DB) error {
			for _, table := range wrongTableNames {
				if err := tx.Migrator().RenameTable(table.Original, table.New); err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			for _, table := range wrongTableNames {
				if err := tx.Migrator().RenameTable(table.New, table.Original); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
