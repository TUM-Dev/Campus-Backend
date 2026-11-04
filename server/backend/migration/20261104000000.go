package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Renaming the rating tables left existing MySQL views pointing to their old
// names. Preserve the view names consumed by the API and repair their sources.
func migrate20261104000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID:      "20261104000000",
		Migrate: refreshRatingViews,
		// The repair is compatible with the preceding application/schema. Rolling
		// back this migration need not restore views that reference missing tables.
		Rollback: func(*gorm.DB) error { return nil },
	}
}

func refreshRatingViews(db *gorm.DB) error {
	for _, statement := range ratingViewStatements {
		if err := db.Exec(statement).Error; err != nil {
			return err
		}
	}
	return nil
}

var ratingViewStatements = []string{
	`CREATE OR REPLACE VIEW cafeteria_rating_statistics AS
 SELECT cafeteriaID, AVG(points) AS average, MIN(points) AS min,
 MAX(points) AS max, STD(points) AS std
 FROM canteen_ratings GROUP BY cafeteriaID`,
	`CREATE OR REPLACE VIEW dish_rating_statistics AS
 SELECT d.cafeteriaID, r.dishID, AVG(r.points) AS average,
 MIN(r.points) AS min, MAX(r.points) AS max, STD(r.points) AS std
 FROM dishes d JOIN dish_ratings r ON r.dishID = d.dish
 GROUP BY d.cafeteriaID, r.dishID`,
	`CREATE OR REPLACE VIEW dish_rating_tag_statistics AS
 SELECT d.cafeteriaID, t.tagID, r.dishID, AVG(t.points) AS average,
 MIN(t.points) AS min, MAX(t.points) AS max, STD(t.points) AS std
 FROM dishes d JOIN dish_ratings r ON r.dishID = d.dish
 JOIN dish_rating_tags t ON t.parentRating = r.dishRating
 GROUP BY d.cafeteriaID, t.tagID, r.dishID`,
	`CREATE OR REPLACE VIEW cafeteria_rating_tag_statistics AS
 SELECT r.cafeteriaID, t.tagID, AVG(t.points) AS average,
 MIN(t.points) AS min, MAX(t.points) AS max, STD(t.points) AS std
 FROM canteen_ratings r JOIN canteen_rating_tags t
 ON t.correspondingRating = r.cafeteriaRating
 GROUP BY r.cafeteriaID, t.tagID`,
	`CREATE OR REPLACE VIEW dish_name_tag_statistics AS
 SELECT d.cafeteriaID, t.tagNameID AS tagID, AVG(t.points) AS average,
 MIN(t.points) AS min, MAX(t.points) AS max, STD(t.points) AS std
 FROM dishes d JOIN dish_ratings r ON r.dishID = d.dish
 JOIN dish_name_tags t ON t.correspondingRating = r.dishRating
 GROUP BY d.cafeteriaID, t.tagNameID`,
}
