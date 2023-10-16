package utils

import (
	"os"

	"github.com/TUM-Dev/Campus-Backend/server/backend/migration"
	gormlogrus "github.com/onrik/gorm-logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SetupDB connects to the database and migrates it if necessary
func SetupDB() *gorm.DB {
	dbHost := os.Getenv("DB_DSN")
	if dbHost == "" {
		log.Fatal("Failed to start! The 'DB_DSN' environment variable is not defined. Take a look at the README.md for more details.")
	}

	log.Info("Connecting to dsn")
	db, err := gorm.Open(mysql.Open(dbHost), &gorm.Config{Logger: gormlogrus.New()})
	if err != nil {
		log.WithError(err).Fatal("failed to connect database")
	}

	// Migrate the schema
	// currently not activated as
	if err := migration.New(db, false).Migrate(); err != nil {
		log.WithError(err).Fatal("Failed to migrate database")
	}
	return db
}
