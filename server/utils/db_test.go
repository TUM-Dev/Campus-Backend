package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAutoMigration(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	autoContainer := setupMySQLTestContainer(ctx, t)
	autoMappedPort, err := autoContainer.MappedPort(ctx, "3306/tcp")
	require.NoError(t, err)
	auto := connectToDbAndMigrate(autoMappedPort, t, true)
	require.NoError(t, auto.Exec("SELECT VERSION()").Error)
}

func TestManualMigration(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	// manual migrations
	manualContainer := setupMySQLTestContainer(ctx, t)
	manualMappedPort, err := manualContainer.MappedPort(ctx, "3306/tcp")
	require.NoError(t, err)
	manual := connectToDbAndMigrate(manualMappedPort, t, false)
	require.NoError(t, manual.Exec("SELECT VERSION()").Error)
}
