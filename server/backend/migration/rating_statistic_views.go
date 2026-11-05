package migration

import (
	"fmt"

	"gorm.io/gorm"
)

// ratingStatisticViews are the views that aggregate canteen and dish ratings.
// They are defined once here, so that manualMigrate and autoMigrate produce the same views.
//
// The view names keep their legacy cafeteria_* names, because the models model.CanteenRatingStatistic
// and model.CanteenRatingTagsStatistic (formerly Cafeteria*) read from them.
var ratingStatisticViews = []struct {
	Name  string
	Query string
}{
	{"cafeteria_rating_statistics", `
SELECT cafeteriaID, AVG(points) AS average, MIN(points) AS min, MAX(points) AS max, STD(points) AS std
FROM canteen_ratings
GROUP BY cafeteriaID
ORDER BY COUNT(*) DESC, average DESC`},
	{"dish_rating_statistics", `
SELECT d.cafeteriaID AS cafeteriaID,
       r.dishID      AS dishID,
       AVG(r.points) AS average,
       MAX(r.points) AS max,
       MIN(r.points) AS min,
       STD(r.points) AS std
FROM dishes d
JOIN dish_ratings r ON r.dishID = d.dish
GROUP BY d.cafeteriaID, r.dishID
ORDER BY COUNT(0) DESC, AVG(r.points) DESC`},
	{"dish_rating_tag_statistics", `
SELECT d.cafeteriaID    AS cafeteriaID,
       mrt.tagID        AS tagID,
       mr.dishID        AS dishID,
       AVG(mrt.points)  AS average,
       MAX(mrt.points)  AS max,
       MIN(mrt.points)  AS min,
       STD(mrt.points)  AS std
FROM dishes d
JOIN dish_ratings mr ON d.dish = mr.dishID
JOIN dish_rating_tags mrt ON mr.dishRating = mrt.parentRating
GROUP BY d.cafeteriaID, mrt.tagID, mr.dishID`},
	{"cafeteria_rating_tag_statistics", `
SELECT cr.cafeteriaID AS cafeteriaID, crt.tagID AS tagID, AVG(crt.points) AS average, MAX(crt.points) AS max, MIN(crt.points) AS min, STD(crt.points) AS std
FROM canteen_ratings cr
JOIN canteen_rating_tags crt ON cr.cafeteriaRating = crt.correspondingRating
GROUP BY cr.cafeteriaID, crt.tagID`},
	{"dish_name_tag_statistics", `
SELECT d.cafeteriaID   AS cafeteriaID,
       mnt.tagNameID   AS tagID,
       AVG(mnt.points) AS average,
       MAX(mnt.points) AS max,
       MIN(mnt.points) AS min,
       STD(mnt.points) AS std
FROM dishes d
JOIN dish_ratings mr ON d.dish = mr.dishID
JOIN dish_name_tags mnt ON mr.dishRating = mnt.correspondingRating
GROUP BY d.cafeteriaID, mnt.tagNameID`},
}

// createRatingStatisticViews creates or replaces all rating statistic views
func createRatingStatisticViews(tx *gorm.DB) error {
	for _, view := range ratingStatisticViews {
		if err := tx.Exec(fmt.Sprintf("CREATE OR REPLACE VIEW `%s` AS %s", view.Name, view.Query)).Error; err != nil {
			return fmt.Errorf("creating view %s: %w", view.Name, err)
		}
	}
	return nil
}
