// Package migration contains functions related to database changes and executes them
package migration

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
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

func manualMigrate(db *gorm.DB) error {
	gormigrateOptions := &gormigrate.Options{
		TableName:                 gormigrate.DefaultOptions.TableName,
		IDColumnName:              gormigrate.DefaultOptions.IDColumnName,
		IDColumnSize:              gormigrate.DefaultOptions.IDColumnSize,
		UseTransaction:            true,
		ValidateUnknownMigrations: true,
	}
	mig := gormigrate.New(db, gormigrateOptions, []*gormigrate.Migration{
		migrate20200000000000(),
		migrate20210709193000(),
		migrate20220126230000(),
		migrate20220713000000(),
		migrate20221119131300(),
		migrate20221210000000(),
		migrate20230825000000(),
		migrate20230904000000(),
		migrate20230530000000(),
		migrate20230904100000(),
		migrate20230826000000(),
		migrate20231003000000(),
	})
	return mig.Migrate()
}

// Migrate starts the migration either by using AutoMigrate in development environments or manually in prod
func Migrate(db *gorm.DB, shouldAutoMigrate bool) error {
	log.WithField("shouldAutoMigrate", shouldAutoMigrate).Debug("starting migration")
	if shouldAutoMigrate {
		return autoMigrate(db)
	} else {
		return manualMigrate(db)
	}
}
