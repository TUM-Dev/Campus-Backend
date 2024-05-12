package utils

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"strings"
	"testing"
	"time"

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
	if err := migration.Migrate(db, os.Getenv("CI_AUTO_MIGRATION") == "true"); err != nil {
		log.WithError(err).Fatal("Failed to migrate database")
	}

	if os.Getenv("CI_EXIT_AFTER_MIGRATION") == "true" {
		log.Info("Exiting after migration")
		os.Exit(0)
	}
	return db
}

type testContainerLogger struct {
	t *testing.T
}

func (tcl testContainerLogger) Printf(format string, v ...interface{}) {
	line := strings.TrimSpace(fmt.Sprintf(format, v...))
	tcl.t.Log(line)
}

func (tcl testContainerLogger) Accept(log testcontainers.Log) {
	line := strings.TrimSpace(string(log.Content))
	if len(line) == 0 {
		return // empty lines are just junk..
	}
	tcl.t.Logf("[%s,testcontainer] %s", log.LogType, line)
}

func SetupTestContainer(ctx context.Context, t *testing.T) *gorm.DB {
	logger := testContainerLogger{t}
	// create a container
	req := testcontainers.ContainerRequest{
		Image: "mysql:8",
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "super_secret_passw0rd",
			"MYSQL_DATABASE":      "campus_db",
		},
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{logger},
		},
		WaitingFor: &wait.LogStrategy{
			Log:          "mysqld: ready for connections",
			IsRegexp:     false,
			Occurrence:   2, // why does it do a dance with a temporary server???
			PollInterval: 100 * time.Millisecond,
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           logger,
	})
	require.NoError(t, err)
	require.True(t, container.IsRunning())
	t.Cleanup(func() {
		require.NoError(t, container.Terminate(ctx))
	})
	// connect to gorm instance
	mappedPort, err := container.MappedPort(ctx, "3306/tcp")
	require.NoError(t, err)
	dsn := fmt.Sprintf("root:super_secret_passw0rd@tcp(localhost:%d)/campus_db?charset=utf8mb4&parseTime=True&loc=Local", mappedPort.Int())
	t.Log("connecting to " + dsn)
	db, err := gorm.Open(mysql.Open(dsn))
	require.NoError(t, err)
	// we run the auto-migrations because they are WAY faster
	require.NoError(t, migration.Migrate(db, true))
	return db
}
