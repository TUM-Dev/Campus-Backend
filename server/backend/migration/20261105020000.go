package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20261105020000
// Retains the ID for databases that already recorded it. The student council
// tables are owned by SQL migration 20261105185936, so this must not change schema.
func migrate20261105020000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20261105020000",
		Migrate: func(*gorm.DB) error {
			return nil
		},
		Rollback: func(*gorm.DB) error {
			return nil
		},
	}
}
