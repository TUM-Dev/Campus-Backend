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
				&cafeteria_rating_models.Dish{},
				&cafeteria_rating_models.DishNameTagOption{},
				&cafeteria_rating_models.DishNameTagOptionIncluded{},
				&cafeteria_rating_models.DishNameTagOptionExcluded{},
				&cafeteria_rating_models.DishNameTag{},
				&cafeteria_rating_models.DishNameTagAverage{},
				&cafeteria_rating_models.DishRating{},
				&cafeteria_rating_models.DishRatingAverage{},
				&cafeteria_rating_models.DishRatingTag{},
				&cafeteria_rating_models.DishRatingTagAverage{},
				&cafeteria_rating_models.DishRatingTagOption{},
				&cafeteria_rating_models.DishToDishNameTag{},
			); err != nil {
				return err
			}
			return nil
		},
	}
}
