package utils

import (
	"testing"

	"github.com/TUM-Dev/Campus-Backend/server/backend/migration"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var ratingStatisticViews = []string{
	"cafeteria_rating_statistics",
	"cafeteria_rating_tag_statistics",
	"dish_rating_statistics",
	"dish_rating_tag_statistics",
	"dish_name_tag_statistics",
}

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

func TestStudentCouncilMigration(t *testing.T) {
	t.Parallel()
	for name, migrate := range map[string]func(*gorm.DB) error{
		"baseline": migration.MigrateFromBaseline, "production": migration.Migrate,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			db := newTestDatabase(t, migrate)
			require.True(t, db.Migrator().HasTable(&model.StudentCouncilCollection{}))
			require.True(t, db.Migrator().HasTable(&model.StudentCouncil{}))

			// Go migrations must not undo enum values added by later SQL migrations.
			var column struct{ Type string }
			require.NoError(t, db.Raw("SHOW COLUMNS FROM crontab LIKE 'type'").Scan(&column).Error)
			require.Contains(t, column.Type, "'tumonlineStatus'")

			if name == "production" {
				require.NoError(t, db.Exec("INSERT INTO student_council_collections (name) VALUES ('preserved')").Error)
				// Reproduce a deployment where SQL succeeded but the obsolete Go migration failed.
				require.NoError(t, db.Exec("DELETE FROM migrations WHERE id = '20261105020000'").Error)
				require.NoError(t, migration.Migrate(db))
				// MySQL DDL can survive an interrupted migration without its ID being recorded.
				require.NoError(t, db.Exec("DELETE FROM migrations WHERE id IN ('20261105185936', '20261106000000')").Error)
				require.NoError(t, db.Exec("DELETE FROM crontab WHERE type = 'tumonlineStatus'").Error)
				require.NoError(t, migration.Migrate(db))
				var count int64
				require.NoError(t, db.Table("student_council_collections").Where("name = ?", "preserved").Count(&count).Error)
				require.EqualValues(t, 1, count)
			}
		})
	}
}

// TestRatingStatisticViewsQueryable makes sure that the rating statistic views reference existing tables.
// MySQL does not update views when a table they reference is renamed. A broken view only fails when it is queried.
func TestRatingStatisticViewsQueryable(t *testing.T) {
	t.Parallel()
	for name, fromBaseline := range map[string]bool{"baseline": true, "production": false} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			migrate := migration.Migrate
			if fromBaseline {
				migrate = migration.MigrateFromBaseline
			}
			db := newTestDatabase(t, migrate)
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
