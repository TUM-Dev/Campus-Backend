package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// CafeteriaRatingAverage stores all precomputed values for the cafeteria ratings
type CafeteriaRatingAverage struct {
	CafeteriaRatingAverage int64   `gorm:"primary_key;AUTO_INCREMENT;column:cafeteriaRatingAverage;type:int;not null;"`
	CafeteriaID            int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;"`
	Average                float64 `gorm:"column:average;type:float;not null;"`
	Min                    int32   `gorm:"column:min;type:int;not null;"`
	Max                    int32   `gorm:"column:max;type:int;not null;"`
	Std                    float64 `gorm:"column:std;type:float;not null;"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRatingAverage) TableName() string {
	return "cafeteria_rating_average"
}

// DishRatingAverage stores all precomputed values for the cafeteria ratings
type DishRatingAverage struct {
	DishRatingAverage int64   `gorm:"primary_key;AUTO_INCREMENT;column:dishRatingAverage;type:int;not null;"`
	CafeteriaID       int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;"`
	DishID            int64   `gorm:"column:dishID;foreignKey:dish;type:int;not null;"`
	Average           float64 `gorm:"column:average;type:float;not null;"`
	Min               int32   `gorm:"column:min;type:int;not null;"`
	Max               int32   `gorm:"column:max;type:int;not null;"`
	Std               float64 `gorm:"column:std;type:float;not null;"`
}

// TableName sets the insert table name for this struct type
func (n *DishRatingAverage) TableName() string {
	return "dish_rating_average"
}

// migrate20231015000000
// migrates the static data for the canteen rating system and adds the necessary cronjob entries
func (m TumDBMigrator) migrate20231015000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20231015000000",
		Migrate: func(tx *gorm.DB) error {
			tables := []string{"cafeteria_rating_average", "dish_rating_average"}
			for _, table := range tables {
				if err := tx.Migrator().DropTable(table); err != nil {
					return err
				}
			}
			if err := tx.Exec(`create view cafeteria_rating_statistics as
SELECT cafeteriaID, Avg(points) AS average, MIN(points) AS min, Max(points) AS max, STD(points) AS std
FROM cafeteria_rating
GROUP BY cafeteriaID
ORDER BY COUNT(cafeteriaID) DESC, average DESC`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`create view dish_rating_statistics as
SELECT cafeteriaID, dishID, AVG(points) as average, MAX(points) as max, MIN(points) as min, STD(points) as std
FROM dish_rating
GROUP BY cafeteriaID,dishID
ORDER BY COUNT(*) DESC, average DESC`).Error; err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			createdViews := []string{"cafeteria_rating_statistics", "dish_rating_statistics"}
			for _, view := range createdViews {
				if err := tx.Exec("DROP VIEW IF EXISTS " + view).Error; err != nil {
					return err
				}
			}
			return tx.AutoMigrate(&CafeteriaRatingAverage{}, &DishRatingAverage{})
		},
	}
}
