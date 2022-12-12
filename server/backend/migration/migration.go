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
			&model.Files{},
			&model.NewsSource{},
			&model.NewsAlert{},
			&model.News{},
			&model.Cafeteria{},
			&model.CafeteriaRating{},
			&model.CafeteriaRatingTag{},
			&model.CafeteriaRatingTagOption{},
			&model.CafeteriaRatingTagsAverage{},
			&model.Dish{},
			&model.DishNameTag{},
			&model.DishNameTagAverage{},
			&model.DishNameTagOption{},
			&model.DishNameTagOptionExcluded{},
			&model.DishNameTagOptionIncluded{},
			&model.DishRating{},
			&model.DishRatingAverage{},
			&model.DishRatingTag{},
			&model.DishRatingTagAverage{},
			&model.DishRatingTagOption{},
			&model.DishToDishNameTag{},
			&model.DishesOfTheWeek{},
			&model.CanteenHomometer{},
		)
		return err
	}
	log.Info("Using manual migration")
	mig := gormigrate.New(m.database, gormigrate.DefaultOptions, []*gormigrate.Migration{
		m.migrate20210709193000(),
		m.migrate20220126230000(),
		m.migrate20220713000000(),
		m.migrate20221210000000(),
	})
	err := mig.Migrate()
	return err

}
