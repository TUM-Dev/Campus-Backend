package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAutoMigration(t *testing.T) {
	t.Parallel()
	auto := newTestDatabase(t, true)
	require.NoError(t, auto.Exec("SELECT VERSION()").Error)
}

func TestManualMigration(t *testing.T) {
	t.Parallel()
	manual := newTestDatabase(t, false)
	require.NoError(t, manual.Exec("SELECT VERSION()").Error)
}
