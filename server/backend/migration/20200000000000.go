package migration

import (
	_ "embed"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

//go:embed static_data/source-schema.sql
var sourceSchema string

// migrate20200000000000
// adds the source schema
func migrate20200000000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20200000000000",
		Migrate: func(tx *gorm.DB) error {
			tx = tx.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})
			for _, line := range strings.Split(sourceSchema, ";") {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				if err := tx.Exec(line).Error; err != nil {
					log.WithError(err).WithField("line", line).Error("failed to execute line")
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			re := regexp.MustCompile(`create table if not exists (?P<table_name>\S+)`)
			tables := re.FindAllStringSubmatch(sourceSchema, -1)
			for _, table := range tables {
				if err := tx.Migrator().DropTable(table[re.SubexpIndex("table_name")]); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
