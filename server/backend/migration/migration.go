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
		&model.Cafeteria{},
		&model.CafeteriaRating{},
		&model.CafeteriaRatingAverage{},
		&model.CafeteriaRatingTag{},
		&model.CafeteriaRatingTagAverage{},
		&model.CafeteriaRatingTagOption{},
		&model.CanteenHeadCount{},
		&model.Crontab{},
		&model.Device{},
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
		&model.PublishedExamResult{},
		&model.Feedback{},
		&model.File{},
		&model.IOSDevice{},
		//&model.IOSDeviceLastUpdated{},
		&model.IOSDeviceRequestLog{},
		&model.IOSDevicesActivityReset{},
		&model.IOSGrade{},
		//&model.IOSRemoteNotification{},
		&model.IOSScheduledUpdateLog{},
		&model.IOSSchedulingPriority{},
		&model.Kino{},
		&model.NewExamResultsSubscriber{},
		&model.News{},
		&model.NewsAlert{},
		&model.NewsSource{},
		&model.Notification{},
		&model.NotificationConfirmation{},
		&model.NotificationType{},
		&model.UpdateNote{},
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
