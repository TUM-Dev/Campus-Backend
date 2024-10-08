// Package migration contains functions related to database changes and executes them
package migration

import (
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.Canteen{},
		&model.CanteenRating{},
		&model.CanteenRatingTag{},
		&model.CanteenRatingTagOption{},
		&model.CanteenHeadCount{},
		&model.Crontab{},
		&model.Device{},
		&model.Dish{},
		&model.DishNameTag{},
		&model.DishNameTagOption{},
		&model.ExcludedDishNameTagOption{},
		&model.IncludedDishNameTagOption{},
		&model.DishRating{},
		&model.DishRatingTag{},
		&model.DishRatingTagOption{},
		&model.DishToDishNameTag{},
		&model.DishesOfTheWeek{},
		&model.Feedback{},
		&model.File{},
		&model.Movie{},
		&model.NewExamResultsSubscriber{},
		&model.News{},
		&model.NewsAlert{},
		&model.NewsSource{},
		&model.Notification{},
		&model.NotificationConfirmation{},
		&model.NotificationType{},
		&model.UpdateNote{},
		&model.StudentClub{},
		&model.StudentClubCollection{},
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
	migrations := []*gormigrate.Migration{
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
		migrate20231015000000(),
		migrate20231023000000(),
		migrate20240101000000(),
		migrate20240102000000(),
		migrate20240103000000(),
		migrate20240207000000(),
		migrate20240311000000(),
		migrate20240312000000(),
		migrate20240316000000(),
		migrate20240317000000(),
		migrate20240318000000(),
		migrate20240319000000(),
		migrate20240320000000(),
		migrate20240327000000(),
		migrate20240402000000(),
		migrate20240405000000(),
		migrate20240510000000(),
		migrate20240511000000(),
		migrate20240512000000(),
		migrate20240706000000(),
		migrate20240824000000(),
		migrate20241023000000(),
	}
	return gormigrate.New(db, gormigrateOptions, migrations).Migrate()
}

// Migrate starts the migration either by using AutoMigrate in development environments or manually in prod
func Migrate(db *gorm.DB, shouldAutoMigrate bool) error {
	// without this set, the `COLLATE` will not be set correctly
	db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")

	log.WithField("shouldAutoMigrate", shouldAutoMigrate).Info("starting migration")
	start := time.Now()
	var err error
	if shouldAutoMigrate {
		err = autoMigrate(db)
	} else {
		err = manualMigrate(db)
	}
	log.WithField("elapsed", time.Since(start)).Info("migration done")
	return err
}
