package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20220126230000
// adds a fulltext index to the roomfinder_rooms table
func (m TumDBMigrator) migrate20220126230000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20220126230000",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec("CREATE FULLTEXT INDEX `search_index` ON `roomfinder_rooms` (`info`, `address`, `room_code`)").Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec("DROP INDEX search_index ON roomfinder_rooms").Error
		},
	}
}
