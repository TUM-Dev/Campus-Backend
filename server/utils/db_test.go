package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestAutoMigration(t *testing.T) {
	t.Parallel()
	auto := newTestDatabase(t, true)
	require.NoError(t, auto.Exec("SELECT VERSION()").Error)
	assertRatingViews(t, auto)
}

func TestManualMigration(t *testing.T) {
	t.Parallel()
	manual := newTestDatabase(t, false)
	require.NoError(t, manual.Exec("SELECT VERSION()").Error)
	assertRatingViews(t, manual)
}

func assertRatingViews(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, view := range []string{"cafeteria_rating_statistics", "dish_rating_statistics", "dish_rating_tag_statistics", "cafeteria_rating_tag_statistics", "dish_name_tag_statistics"} {
		rows, err := db.Raw("SELECT * FROM " + view + " LIMIT 0").Rows()
		require.NoError(t, err, "view %s must resolve after migration", view)
		require.NoError(t, rows.Close())
	}
}
