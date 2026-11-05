package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ratingStatisticViews are the views that back the rating endpoints
var ratingStatisticViews = []string{
	"cafeteria_rating_statistics",
	"cafeteria_rating_tag_statistics",
	"dish_rating_statistics",
	"dish_rating_tag_statistics",
	"dish_name_tag_statistics",
}

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

// TestRatingStatisticViewsQueryable makes sure that the rating statistic views reference existing tables.
// MySQL does not update views when a table they reference is renamed. A broken view only fails when it is queried.
func TestRatingStatisticViewsQueryable(t *testing.T) {
	t.Parallel()
	for name, autoMigrate := range map[string]bool{"auto": true, "manual": false} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			db := newTestDatabase(t, autoMigrate)
			for _, view := range ratingStatisticViews {
				var rows []map[string]any
				require.NoError(t, db.Table(view).Limit(1).Find(&rows).Error, view)
			}

			// the views must aggregate the data of the renamed tables
			require.NoError(t, db.Exec("INSERT INTO canteens (cafeteria, name, address, latitude, longitude) VALUES (1, 'MENSA_GARCHING', 'Boltzmannstr. 19', 48.26, 11.67)").Error)
			require.NoError(t, db.Exec("INSERT INTO canteen_ratings (points, comment, cafeteriaID) VALUES (2, '', 1), (4, '', 1)").Error)
			require.NoError(t, db.Exec("INSERT INTO dishes (dish, name, type, cafeteriaID) VALUES (1, 'Pasta', 'Pasta', 1)").Error)
			require.NoError(t, db.Exec("INSERT INTO dish_ratings (points, dishID, comment) VALUES (1, 1, ''), (5, 1, '')").Error)

			var canteenAvg, dishAvg float64
			require.NoError(t, db.Raw("SELECT average FROM cafeteria_rating_statistics WHERE cafeteriaID = 1").Scan(&canteenAvg).Error)
			require.NoError(t, db.Raw("SELECT average FROM dish_rating_statistics WHERE cafeteriaID = 1 AND dishID = 1").Scan(&dishAvg).Error)
			require.InDelta(t, 3.0, canteenAvg, 0.001)
			require.InDelta(t, 3.0, dishAvg, 0.001)
		})
	}
}
