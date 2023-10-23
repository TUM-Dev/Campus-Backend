package migration

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Cafeteria stores all Available cafeterias in the format of the eat-api
type Cafeteria struct {
	Cafeteria int64   `gorm:"primary_key;AUTO_INCREMENT;column:cafeteria;type:int;not null;" json:"canteen" `
	Name      string  `gorm:"column:name;type:mediumtext;not null;" json:"name" `
	Address   string  `gorm:"column:address;type:text;not null;" json:"address" `
	Latitude  float32 `gorm:"column:latitude;type:float;not null;" json:"latitude" `
	Longitude float32 `gorm:"column:longitude;type:float;not null;" json:"longitude"`
}

// TableName sets the insert table name for this struct type
func (n *Cafeteria) TableName() string {
	return "cafeteria"
}

//migrate20210709193000

func migrate20220713000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20220713000000",
		Migrate: func(tx *gorm.DB) error {

			if err := tx.AutoMigrate(
				&Cafeteria{},
				&model.CanteenRating{},
				&model.CafeteriaRatingAverage{},
				&model.CanteenRatingTag{},
				&model.CafeteriaRatingTagAverage{},
				&model.CanteenRatingTagOption{},
				&model.Dish{},
				&model.DishesOfTheWeek{},
				&model.DishNameTagOption{},
				&model.IncludedDishNameTagOption{},
				&model.ExcludedDishNameTagOption{},
				&model.DishNameTag{},
				&model.DishNameTagAverage{},
				&model.DishRating{},
				&model.DishRatingAverage{},
				&model.DishRatingTag{},
				&model.DishRatingTagAverage{},
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
