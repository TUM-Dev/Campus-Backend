package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

//migrate20220126230000
//adds a fulltext index to the roomfinder_rooms table
func (m TumDBMigrator) migrate20220126230000() *gormigrate.Migration {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	m.database.Logger = newLogger
	return &gormigrate.Migration{
		ID: "20220126230000",
		Migrate: func(tx *gorm.DB) error {
			return tx.Debug().Exec("CREATE FULLTEXT INDEX `search_index` ON `roomfinder_rooms` (`info`, `address`, `room_code`)").Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec("DROP INDEX search_index ON roomfinder_rooms").Error
		},
	}
}
