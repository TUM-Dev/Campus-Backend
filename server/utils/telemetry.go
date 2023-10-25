package utils

import (
	"os"

	"github.com/TUM-Dev/Campus-Backend/server/env"
	"github.com/getsentry/sentry-go"
	"github.com/makasim/sentryhook"
	log "github.com/sirupsen/logrus"
)

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
