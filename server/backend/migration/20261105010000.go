package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20261105010000
//   - recreates the rating statistic views. They were created in 20231015000000 against the old table names
//     (cafeteria_rating, dish, ...). 20240824000000 renamed those tables, and MySQL does not update views
//     on RENAME TABLE, so every query against these views failed with
//     "Error 1356: View ... references invalid table(s)".
func migrate20261105010000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID:      "20261105010000",
		Migrate: createRatingStatisticViews,
		Rollback: func(tx *gorm.DB) error {
			// The previous view definitions referenced tables that no longer exist and could not be queried.
			// Restoring them would only break the views again, so the rollback keeps the working views.
			return nil
		},
	}
}
