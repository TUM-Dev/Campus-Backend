// Package migration contains functions related to database changes and executes them
package migration

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// TumDBMigrator contains a reference to our database
type TumDBMigrator struct {
	database          *gorm.DB
	shouldAutoMigrate bool
}

// New creates a new TumDBMigrator with a database
func New(db *gorm.DB, shouldAutoMigrate bool) TumDBMigrator {
	return TumDBMigrator{database: db, shouldAutoMigrate: shouldAutoMigrate}
}

// Migrate starts the migration either by using AutoMigrate in development environments or manually in prod
func (m TumDBMigrator) Migrate() error {
	if m.shouldAutoMigrate {
		log.Info("Using automigration")
		err := m.database.AutoMigrate(
			&model.TopNews{},
			&model.Crontab{},
			&model.File{},
			&model.NewsSource{},
			&model.NewsAlert{},
			&model.News{},
			&model.CanteenHeadCount{},
		)
		return err
	}
	log.Info("Using manual migration")
	gormigrateOptions := &gormigrate.Options{
		TableName:                 gormigrate.DefaultOptions.TableName,
		IDColumnName:              gormigrate.DefaultOptions.IDColumnName,
		IDColumnSize:              gormigrate.DefaultOptions.IDColumnSize,
		UseTransaction:            true,
		ValidateUnknownMigrations: true,
	}
	mig := gormigrate.New(m.database, gormigrateOptions, []*gormigrate.Migration{
		m.migrate20210709193000(),
		m.migrate20220126230000(),
		m.migrate20220713000000(),
		m.migrate20221119131300(),
		m.migrate20221210000000(),
		m.migrate20230825000000(),
		m.migrate20230904000000(),
		m.migrate20230530000000(),
		m.migrate20230904100000(),
		m.migrate20230826000000(),
	})
	err := mig.Migrate()
	return err

}
