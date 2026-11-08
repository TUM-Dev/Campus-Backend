package migration

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGoMigrationsFrozen(t *testing.T) {
	var ids []string
	for _, m := range append(legacyMigrations(), postBaselineMigrations()...) {
		ids = append(ids, m.ID)
	}
	// New migrations belong in sql/, not in the frozen Go migration registry.
	require.Equal(t, []string{
		"20200000000000", "20210709193000", "20220126230000", "20220713000000",
		"20221119131300", "20221210000000", "20230825000000", "20230904000000",
		"20230530000000", "20230904100000", "20230826000000", "20231003000000",
		"20231015000000", "20231023000000", "20240101000000", "20240102000000",
		"20240103000000", "20240207000000", "20240311000000", "20240312000000",
		"20240316000000", "20240317000000", "20240318000000", "20240319000000",
		"20240320000000", "20240327000000", "20240402000000", "20240405000000",
		"20240510000000", "20240511000000", "20240512000000", "20240706000000",
		"20240824000000", "20241023000000", "20250412000000",
		"20261105000000", "20261105010000", "20261105020000",
	}, ids)
}

func TestStudentCouncilGoMigrationIsNoOp(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{})
	require.NoError(t, err)

	m := migrate20261105020000()
	require.Equal(t, "20261105020000", m.ID)
	require.NoError(t, m.Migrate(db))
	require.NoError(t, m.Rollback(db))
	require.NoError(t, mock.ExpectationsWereMet())
}

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
	for _, m := range postBaselineMigrations() {
		require.Greater(t, m.ID, baselineVersion)
		require.Less(t, m.ID, files[1].version, "Go migrations must precede all post-baseline SQL migrations")
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
