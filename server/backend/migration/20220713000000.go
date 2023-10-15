package migration

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// InitialCafeteriaRatingAverage stores all precomputed values for the cafeteria ratings
type InitialCafeteriaRatingAverage struct {
	CafeteriaRatingAverage int64   `gorm:"primary_key;AUTO_INCREMENT;column:cafeteriaRatingAverage;type:int;not null;"`
	CafeteriaID            int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;"`
	Average                float64 `gorm:"column:average;type:float;not null;"`
	Min                    int32   `gorm:"column:min;type:int;not null;"`
	Max                    int32   `gorm:"column:max;type:int;not null;"`
	Std                    float64 `gorm:"column:std;type:float;not null;"`
}

// TableName sets the insert table name for this struct type
func (n *InitialCafeteriaRatingAverage) TableName() string {
	return "cafeteria_rating_average"
}

// InitialDishRatingAverage stores all precomputed values for the cafeteria ratings
type InitialDishRatingAverage struct {
	DishRatingAverage int64   `gorm:"primary_key;AUTO_INCREMENT;column:dishRatingAverage;type:int;not null;"`
	CafeteriaID       int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;"`
	DishID            int64   `gorm:"column:dishID;foreignKey:dish;type:int;not null;"`
	Average           float64 `gorm:"column:average;type:float;not null;"`
	Min               int32   `gorm:"column:min;type:int;not null;"`
	Max               int32   `gorm:"column:max;type:int;not null;"`
	Std               float64 `gorm:"column:std;type:float;not null;"`
}

// TableName sets the insert table name for this struct type
func (n *InitialDishRatingAverage) TableName() string {
	return "dish_rating_average"
}

// InitialDishRatingTagAverage stores all precomputed values for the cafeteria ratings
type InitialDishRatingTagAverage struct {
	DishRatingTagsAverage int64   `gorm:"primary_key;AUTO_INCREMENT;column:dishRatingTagsAverage;type:int;not null;"`
	CafeteriaID           int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;"`
	TagID                 int64   `gorm:"column:tagID;foreignKey:tagID;type:int;not null;"`
	DishID                int64   `gorm:"column:dishID;foreignKey:dishID;type:int;not null;"`
	Average               float32 `gorm:"column:average;type:float;not null;"`
	Min                   int8    `gorm:"column:min;type:int;not null;"`
	Max                   int8    `gorm:"column:max;type:int;not null;"`
	Std                   float32 `gorm:"column:std;type:float;not null;"`
}

// TableName sets the insert table name for this struct type
func (n *InitialDishRatingTagAverage) TableName() string {
	return "dish_rating_tag_average"
}

// InitialCafeteriaRatingTagsAverage stores all precomputed values for the cafeteria ratings
type InitialCafeteriaRatingTagsAverage struct {
	CafeteriaRatingTagsAverage int64   `gorm:"primary_key;AUTO_INCREMENT;column:cafeteriaRatingTagsAverage;type:int;not null;" json:"canteenRatingTagsAverage"`
	CafeteriaID                int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"canteenID"`
	TagID                      int64   `gorm:"column:tagID;foreignKey:cafeteriaRatingTagOption;type:int;not null;" json:"tagID"`
	Average                    float32 `gorm:"column:average;type:float;not null;" json:"average"`
	Min                        int8    `gorm:"column:min;type:int;not null;" json:"min"`
	Max                        int8    `gorm:"column:max;type:int;not null;" json:"max"`
	Std                        float32 `gorm:"column:std;type:float;not null;" json:"std"`
}

// TableName sets the insert table name for this struct type
func (n *InitialCafeteriaRatingTagsAverage) TableName() string {
	return "cafeteria_rating_tag_average"
}

// migrate20220713000000
func (m TumDBMigrator) migrate20220713000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20220713000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(
				&model.Cafeteria{},
				&model.CafeteriaRating{},
				&InitialCafeteriaRatingAverage{},
				&model.CafeteriaRatingTag{},
				&InitialCafeteriaRatingTagsAverage{},
				&model.CafeteriaRatingTagOption{},
				&model.Dish{},
				&model.DishesOfTheWeek{},
				&model.DishNameTagOption{},
				&model.DishNameTagOptionIncluded{},
				&model.DishNameTagOptionExcluded{},
				&model.DishNameTag{},
				&model.DishNameTagAverage{},
				&model.DishRating{},
				&InitialDishRatingAverage{},
				&model.DishRatingTag{},
				&InitialDishRatingTagAverage{},
				&model.DishRatingTagOption{},
				&model.DishToDishNameTag{},
			); err != nil {
				return err
			}
			return nil
		},

		Rollback: func(tx *gorm.DB) error {
			res := tx.Delete(&model.Crontab{}, "type = 'dishNameDownload'").Error
			if res != nil {
				return res
			}
			return tx.Delete(&model.Crontab{}, "type = 'averageRatingComputation'").Error
		},
	}
}
