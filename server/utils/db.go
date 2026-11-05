package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/TUM-Dev/Campus-Backend/server/backend/migration"
	gormlogrus "github.com/onrik/gorm-logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

// TestMySQLDSNEnv names an optional environment variable pointing tests at an existing MySQL server
// (e.g. a CI service container), formatted like "root:password@tcp(localhost:3306)/".
// When unset, a single MySQL testcontainer is started lazily per test binary and shared by all tests.
const TestMySQLDSNEnv = "TEST_MYSQL_DSN"

var (
	testServerOnce sync.Once
	testServerCfg  *gomysql.Config
	testServerErr  error
)

// testServerConfig returns the connection config of the shared MySQL server (without a database selected)
func testServerConfig(t *testing.T) *gomysql.Config {
	testServerOnce.Do(func() {
		if dsn := os.Getenv(TestMySQLDSNEnv); dsn != "" {
			testServerCfg, testServerErr = gomysql.ParseDSN(dsn)
			return
		}
		// not tied to a test context: the container outlives individual tests and is removed by ryuk once the test binary exits
		testServerCfg, testServerErr = startMySQLTestContainer(context.Background())
	})
	require.NoError(t, testServerErr)
	cfg := testServerCfg.Clone()
	cfg.DBName = ""
	cfg.ParseTime = true
	cfg.Loc = time.Local
	cfg.Params = map[string]string{"charset": "utf8mb4"}
	return cfg
}

func startMySQLTestContainer(ctx context.Context) (*gomysql.Config, error) {
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mysql:9.0.0",
			ExposedPorts: []string{"3306/tcp"},
			Env:          map[string]string{"MYSQL_ROOT_PASSWORD": "super_secret_passw0rd"},
			// durability is irrelevant for throwaway test data
			Tmpfs: map[string]string{"/var/lib/mysql": "rw"},
			Cmd: []string{
				"--skip-log-bin",
				"--innodb-flush-log-at-trx-commit=0",
				"--performance-schema=OFF",
			},
			WaitingFor: wait.ForAll(
				// the entrypoint first starts a temporary server for initialisation
				wait.ForLog("mysqld: ready for connections").WithOccurrence(2).WithPollInterval(100*time.Millisecond),
				wait.ForListeningPort("3306/tcp"),
			),
		},
		Started: true,
	})
	if err != nil {
		return nil, fmt.Errorf("starting mysql testcontainer: %w", err)
	}
	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := container.MappedPort(ctx, "3306/tcp")
	if err != nil {
		return nil, err
	}
	cfg := gomysql.NewConfig()
	cfg.User = "root"
	cfg.Passwd = "super_secret_passw0rd"
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("%s:%d", host, port.Num())
	return cfg, nil
}

// SetupTestContainer returns a connection to a fresh, auto-migrated database on the shared test MySQL server
func SetupTestContainer(_ context.Context, t *testing.T) *gorm.DB {
	return newTestDatabase(t, true)
}

// newTestDatabase creates an isolated database for this test and migrates it
//
// The option to allow for the auto-migrations is because they are WAY faster and this is an option for testing reasons
func newTestDatabase(t *testing.T, shouldAutoMigrate bool) *gorm.DB {
	cfg := testServerConfig(t)
	admin, err := gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{Logger: logger.Discard})
	require.NoError(t, err)

	name := "test_" + strings.ReplaceAll(rand.Text(), "-", "")[:16]
	require.NoError(t, admin.Exec(fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", name)).Error)

	cfg.DBName = name
	db, err := gorm.Open(mysql.Open(cfg.FormatDSN()))
	require.NoError(t, err)
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
		require.NoError(t, admin.Exec(fmt.Sprintf("DROP DATABASE `%s`", name)).Error)
		if sqlDB, err := admin.DB(); err == nil {
			_ = sqlDB.Close()
		}
	})
	require.NoError(t, migration.Migrate(db, shouldAutoMigrate))
	return db
}
