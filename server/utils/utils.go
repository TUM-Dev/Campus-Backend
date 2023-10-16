package utils

import (
	"os"

	"github.com/TUM-Dev/Campus-Backend/server/backend/migration"
	"github.com/TUM-Dev/Campus-Backend/server/env"
	"github.com/getsentry/sentry-go"
	"github.com/makasim/sentryhook"
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

// SetupTelemetry initializes our telemetry stack
// - sentry to be connected with log
// - logrus to
func SetupTelemetry(Version string) {
	environment := "development"
	log.SetLevel(log.TraceLevel)
	if env.IsProd() {
		log.SetLevel(log.InfoLevel)
		environment = "production"
		log.SetFormatter(&log.JSONFormatter{}) // simpler to query but harder to parse in the console
	}

	if sentryDSN := os.Getenv("SENTRY_DSN"); sentryDSN != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDSN,
			AttachStacktrace: true,
			Release:          Version,
			Dist:             Version, // see https://github.com/getsentry/sentry-react-native/issues/516 why this is equal
			Environment:      environment,
		}); err != nil {
			log.WithError(err).Error("Sentry initialization failed")
		}
		log.AddHook(sentryhook.New([]log.Level{log.PanicLevel, log.FatalLevel, log.ErrorLevel, log.WarnLevel}))
	} else {
		log.Info("continuing without sentry")
	}
}
