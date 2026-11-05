package utils

import (
	"testing"

	"github.com/TUM-Dev/Campus-Backend/server/backend/migration"
	"github.com/stretchr/testify/require"
)

func TestMigrateFromBaseline(t *testing.T) {
	t.Parallel()
	db := newTestDatabase(t, migration.MigrateFromBaseline)
	require.NoError(t, db.Exec("SELECT VERSION()").Error)
}

func TestMigrate(t *testing.T) {
	t.Parallel()
	db := newTestDatabase(t, migration.Migrate)
	require.NoError(t, db.Exec("SELECT VERSION()").Error)
	// running it again must be a no-op
	require.NoError(t, migration.Migrate(db))
}
