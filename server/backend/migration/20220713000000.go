package migration

import (
	"github.com/TUM-Dev/Campus-Backend/model/cafeteria_rating_models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

//migrate20210709193000

func (m TumDBMigrator) migrate20220713000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20220713000000",
		Migrate: func(tx *gorm.DB) error {

			if err := tx.AutoMigrate(
				&cafeteria_rating_models.Cafeteria{},
				&cafeteria_rating_models.CafeteriaRating{},
				&cafeteria_rating_models.CafeteriaRatingAverage{},
				&cafeteria_rating_models.CafeteriaRatingTag{},
				&cafeteria_rating_models.CafeteriaRatingTagsAverage{},
				&cafeteria_rating_models.CafeteriaRatingTagOption{},
				&cafeteria_rating_models.Meal{},
				&cafeteria_rating_models.MealNameTagOption{},
				&cafeteria_rating_models.MealNameTagOptionIncluded{},
				&cafeteria_rating_models.MealNameTagOptionExcluded{},
				&cafeteria_rating_models.MealNameTag{},
				&cafeteria_rating_models.MealNameTagAverage{},
				&cafeteria_rating_models.MealRating{},
				&cafeteria_rating_models.MealRatingAverage{},
				&cafeteria_rating_models.MealRatingTag{},
				&cafeteria_rating_models.MealRatingTagAverage{},
				&cafeteria_rating_models.MealRatingTagOption{},
				&cafeteria_rating_models.MealToMealNameTag{},
			); err != nil {
				return err
			}
			return nil
		},
	}
}
