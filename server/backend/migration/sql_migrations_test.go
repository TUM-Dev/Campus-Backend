package migration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSqlMigrationFiles(t *testing.T) {
	files, err := sqlMigrationFiles()
	require.NoError(t, err)
	require.Equal(t, baselineVersion, files[0].version)
	require.NotEmpty(t, files[0].stmts)
	for i := 1; i < len(files); i++ {
		require.Greater(t, files[i].version, files[i-1].version, "versions must be unique and increasing")
	}
	legacyIDs := map[string]bool{}
	for _, m := range legacyMigrations() {
		legacyIDs[m.ID] = true
	}
	for _, f := range files[1:] {
		require.False(t, legacyIDs[f.version], "%s collides with a legacy migration id", f.name)
	}
}

func TestSplitStatements(t *testing.T) {
	stmts := splitStatements("-- Create \"a\" table\nCREATE TABLE `a` (\n  `id` bigint NOT NULL\n);\n-- Modify \"b\" table\nALTER TABLE `b` ADD COLUMN `x` int NULL;\n")
	require.Equal(t, []string{"CREATE TABLE `a` (\n  `id` bigint NOT NULL\n)", "ALTER TABLE `b` ADD COLUMN `x` int NULL"}, stmts)
}
