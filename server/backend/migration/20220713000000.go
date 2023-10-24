package migration

import (
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// InitialCafeteria stores all Available cafeterias in the format of the eat-api
type InitialCafeteria struct {
	Cafeteria int64   `gorm:"primary_key;AUTO_INCREMENT;column:cafeteria;type:int;not null;" json:"canteen" `
	Name      string  `gorm:"column:name;type:mediumtext;not null;" json:"name" `
	Address   string  `gorm:"column:address;type:text;not null;" json:"address" `
	Latitude  float32 `gorm:"column:latitude;type:float;not null;" json:"latitude" `
	Longitude float32 `gorm:"column:longitude;type:float;not null;" json:"longitude"`
}

// TableName sets the insert table name for this struct type
func (n *InitialCafeteria) TableName() string {
	return "cafeteria"
}

// InitialCafeteriaRating stores all Available cafeterias in the format of the eat-api
type InitialCafeteriaRating struct {
	CafeteriaRating int64     `gorm:"primary_key;AUTO_INCREMENT;column:cafeteriaRating;type:int;not null;" json:"canteenrating"`
	Points          int32     `gorm:"column:points;type:int;not null;" json:"points"`
	Comment         string    `gorm:"column:comment;type:text;" json:"comment" `
	CafeteriaID     int64     `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"canteenID"`
	Timestamp       time.Time `gorm:"column:timestamp;type:timestamp;not null;" json:"timestamp" `
	Image           string    `gorm:"column:image;type:text;" json:"image"`
}

// TableName sets the insert table name for this struct type
func (n *InitialCafeteriaRating) TableName() string {
	return "cafeteria_rating"
}

// InitialCafeteriaRatingTag struct is a row record of the either the dish_tag_rating-table or the cafeteria_rating_tags-table in the database
type InitialCafeteriaRatingTag struct {
	CafeteriaRatingTag  int64 `gorm:"primary_key;AUTO_INCREMENT;column:CafeteriaRatingTag;type:int;not null;" json:"CanteenRatingTag" `
	CorrespondingRating int64 `gorm:"foreignKey:cafeteriaRatingID;column:correspondingRating;type:int;not null;" json:"correspondingRating"`
	Points              int32 `gorm:"column:points;type:int;not null;" json:"points"`
	TagID               int64 `gorm:"foreignKey:cafeteriaRatingTagOption;column:tagID;type:int;not null;" json:"tagID"`
}

// TableName sets the insert table name for this struct type
func (n *InitialCafeteriaRatingTag) TableName() string {
	return "cafeteria_rating_tag"
}

// InitialCafeteriaRatingTagOption stores all available options for tags which can be used to quickly rate cafeterias
type InitialCafeteriaRatingTagOption struct {
	CafeteriaRatingsTagOption int64  `gorm:"primary_key;AUTO_INCREMENT;column:cafeteriaRatingTagOption;type:int;not null;" json:"canteenRatingTagOption"`
	DE                        string `gorm:"column:DE;text;default:de;not null;" json:"DE"`
	EN                        string `gorm:"column:EN;text;default:en;not null;" json:"EN"`
}

// TableName sets the insert table name for this struct type
func (n *InitialCafeteriaRatingTagOption) TableName() string {
	return "cafeteria_rating_tag_option"
}

// InitialDish represents one dish fin a specific cafeteria
type InitialDish struct {
	Dish        int64  `gorm:"primary_key;AUTO_INCREMENT;column:dish;type:int;not null;" json:"dish"`
	Name        string `gorm:"column:name;type:text;not null;" json:"name" `
	Type        string `gorm:"column:type;type:text;not null;" json:"type" `
	CafeteriaID int64  `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"cafeteriaID"`
}

// TableName sets the insert table name for this struct type
func (n *InitialDish) TableName() string {
	return "dish"
}

type InitialDishesOfTheWeek struct {
	DishesOfTheWeek int64 `gorm:"primary_key;AUTO_INCREMENT;column:dishesOfTheWeek;type:int;not null;" json:"dishesOfTheWeek"`
	Year            int32 `gorm:"column:year;type:int;not null;" json:"year"`
	Week            int32 `gorm:"column:week;type:int;not null;" json:"week"`
	Day             int32 `gorm:"column:day;type:int;not null;" json:"day"`
	DishID          int64 `gorm:"column:dishID;foreignKey:dish;type:int;not null;" json:"dishID"`
}

// TableName sets the insert table name for this struct type
func (n *InitialDishesOfTheWeek) TableName() string {
	return "dishes_of_the_week"
}

type InitialDishNameTagOption struct {
	DishNameTagOption int64  `gorm:"primary_key;AUTO_INCREMENT;column:dishNameTagOption;type:int;not null;" json:"dishNameTagOption"`
	DE                string `gorm:"column:DE;type:text;not null;" json:"DE"`
	EN                string `gorm:"column:EN;type:text;not null;" json:"EN"`
}

// TableName sets the insert table name for this struct type
func (n *InitialDishNameTagOption) TableName() string {
	return "dish_name_tag_option"
}

type InitialDishRatingTagOption struct {
	DishRatingTagOption int64  `gorm:"primary_key;AUTO_INCREMENT;column:dishRatingTagOption;type:int;not null;" json:"dishRatingTagOption"`
	DE                  string `gorm:"column:DE;type:text;default:de;not null;" json:"DE"`
	EN                  string `gorm:"column:EN;type:text;default:en;not null;" json:"EN"`
}

// TableName sets the insert table name for this struct type
func (n *InitialDishRatingTagOption) TableName() string {
	return "dish_rating_tag_option"
}

type InitialDishNameTag struct {
	DishNameTag         int64 `gorm:"primary_key;AUTO_INCREMENT;column:DishNameTag;type:int;not null;" json:"DishNameTag"`
	CorrespondingRating int64 `gorm:"foreignKey:dish;column:correspondingRating;type:int;not null;" json:"correspondingRating"`
	Points              int32 `gorm:"column:points;type:int;not null;" json:"points"`
	TagNameID           int64 `gorm:"foreignKey:tagRatingID;column:tagNameID;type:int;not null;" json:"tagnameID"`
}

// TableName sets the insert table name for this struct type
func (n *InitialDishNameTag) TableName() string {
	return "dish_name_tag"
}

type InitialDishNameTagOptionExcluded struct {
	DishNameTagOptionExcluded int64  `gorm:"primary_key;AUTO_INCREMENT;column:dishNameTagOptionExcluded;type:int;not null;" json:"dishNameTagOptionExcluded"`
	NameTagID                 int64  `gorm:"foreignKey:dishNameTagOption;column:nameTagID;type:int;not null;" json:"nameTagID"`
	Expression                string `gorm:"column:expression;type:text;" json:"expression"`
}

// TableName sets the insert table name for this struct type
func (n *InitialDishNameTagOptionExcluded) TableName() string {
	return "dish_name_tag_option_excluded"
}

type InitialDishNameTagOptionIncluded struct {
	DishNameTagOptionIncluded int64  `gorm:"primary_key;AUTO_INCREMENT;column:dishNameTagOptionIncluded;type:int;not null;" json:"dishNameTagOptionIncluded"`
	NameTagID                 int64  `gorm:"foreignKey:dishNameTagOption;column:nameTagID;type:int;not null;" json:"nameTagID"`
	Expression                string `gorm:"column:expression;type:text;" json:"expression"`
}

// TableName sets the insert table name for this struct type
func (n *InitialDishNameTagOptionIncluded) TableName() string {
	return "dish_name_tag_option_included"
}

type InitialDishRating struct {
	DishRating  int64     `gorm:"primary_key;AUTO_INCREMENT;column:dishRating;type:int;not null;" json:"dishRating"`
	Points      int32     `gorm:"column:points;type:int;not null;" json:"points"`
	CafeteriaID int64     `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"cafeteriaID"`
	DishID      int64     `gorm:"column:dishID;foreignKey:dish;type:int;not null;" json:"dishID"`
	Comment     string    `gorm:"column:comment;type:text;" json:"comment"`
	Timestamp   time.Time `gorm:"column:timestamp;type:timestamp;not null;" json:"timestamp"`
	Image       string    `gorm:"column:image;type:text;" json:"image"`
}

// TableName sets the insert table name for this struct type
func (n *InitialDishRating) TableName() string {
	return "dish_rating"
}

type InitialDishRatingTag struct {
	DishRatingTag       int64 `gorm:"primary_key;AUTO_INCREMENT;column:dishRatingTag;type:int;not null;" json:"dishRatingTag"`
	CorrespondingRating int64 `gorm:"foreignKey:cafeteriaRating;column:parentRating;type:int;not null;" json:"parentRating"`
	Points              int32 `gorm:"column:points;type:int;not null;" json:"points"`
	TagID               int64 `gorm:"foreignKey:dishRatingTagOption;column:tagID;type:int;not null;" json:"tagID"`
}

// TableName sets the insert table name for this struct type
func (n *InitialDishRatingTag) TableName() string {
	return "dish_rating_tag"
}

type InitialDishToDishNameTag struct {
	DishToDishNameTag int64 `gorm:"primary_key;AUTO_INCREMENT;column:dishToDishNameTag;type:int;not null;" json:"dishToDishNameTag"`
	DishID            int64 `gorm:"column:dishID;foreignKey:dish;type:int;not null;" json:"dishID"`
	NameTagID         int64 `gorm:"foreignKey:dishNameTagOption;column:nameTagID;type:int;not null;" json:"nameTagID"`
}

// TableName sets the insert table name for this struct type
func (n *InitialDishToDishNameTag) TableName() string {
	return "dish_to_dish_name_tag"
}

// InitialCafeteriaRatingAverage stores all precomputed values for the cafeteria ratings
type InitialCafeteriaRatingAverage struct {
	CafeteriaRatingAverage int64   `gorm:"primary_key;AUTO_INCREMENT;column:cafeteriaRatingAverage;type:int;not null;" json:"canteenRatingAverage" `
	CafeteriaID            int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"canteenID"`
	Average                float64 `gorm:"column:average;type:float;not null;" json:"average" `
	Min                    int32   `gorm:"column:min;type:int;not null;" json:"min"`
	Max                    int32   `gorm:"column:max;type:int;not null;" json:"max"`
	Std                    float64 `gorm:"column:std;type:float;not null;" json:"std"`
}

// TableName sets the insert table name for this struct type
func (n *InitialCafeteriaRatingAverage) TableName() string {
	return "cafeteria_rating_average"
}

// InitialCafeteriaRatingTagAverage stores all precomputed values for the cafeteria ratings
type InitialCafeteriaRatingTagAverage struct {
	CafeteriaRatingTagsAverage int64   `gorm:"primary_key;AUTO_INCREMENT;column:cafeteriaRatingTagsAverage;type:int;not null;" json:"canteenRatingTagsAverage"`
	CafeteriaID                int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"canteenID"`
	TagID                      int64   `gorm:"column:tagID;foreignKey:cafeteriaRatingTagOption;type:int;not null;" json:"tagID"`
	Average                    float32 `gorm:"column:average;type:float;not null;" json:"average"`
	Min                        int8    `gorm:"column:min;type:int;not null;" json:"min"`
	Max                        int8    `gorm:"column:max;type:int;not null;" json:"max"`
	Std                        float32 `gorm:"column:std;type:float;not null;" json:"std"`
}

// TableName sets the insert table name for this struct type
func (n *InitialCafeteriaRatingTagAverage) TableName() string {
	return "cafeteria_rating_tag_average"
}

// InitialDishNameTagAverage stores all precomputed values for the DishName ratings
type InitialDishNameTagAverage struct {
	DishNameTagAverage int64   `gorm:"primary_key;AUTO_INCREMENT;column:dishNameTagAverage;type:int;not null;" json:"dishNameTagAverage" `
	CafeteriaID        int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"cafeteriaID"`
	TagID              int64   `gorm:"column:tagID;foreignKey:DishNameTagOption;type:int;not null;" json:"tagID"`
	Average            float32 `gorm:"column:average;type:float;not null;" json:"average" `
	Min                int8    `gorm:"column:min;type:int;not null;" json:"min"`
	Max                int8    `gorm:"column:max;type:int;not null;" json:"max"`
	Std                float32 `gorm:"column:std;type:float;not null;" json:"std"`
}

// TableName sets the insert table name for this struct type
func (n *InitialDishNameTagAverage) TableName() string {
	return "dish_name_tag_average"
}

// InitialDishRatingAverage stores all precomputed values for the cafeteria ratings
type InitialDishRatingAverage struct {
	DishRatingAverage int64   `gorm:"primary_key;AUTO_INCREMENT;column:dishRatingAverage;type:int;not null;" json:"dishRatingAverage" `
	CafeteriaID       int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"cafeteriaID"`
	DishID            int64   `gorm:"column:dishID;foreignKey:dish;type:int;not null;" json:"dishID"`
	Average           float64 `gorm:"column:average;type:float;not null;" json:"average" `
	Min               int32   `gorm:"column:min;type:int;not null;" json:"min"`
	Max               int32   `gorm:"column:max;type:int;not null;" json:"max"`
	Std               float64 `gorm:"column:std;type:float;not null;" json:"std"`
}

// TableName sets the insert table name for this struct type
func (n *InitialDishRatingAverage) TableName() string {
	return "dish_rating_average"
}

// InitialDishRatingTagAverage stores all precomputed values for the cafeteria ratings
type InitialDishRatingTagAverage struct {
	DishRatingTagsAverage int64   `gorm:"primary_key;AUTO_INCREMENT;column:dishRatingTagsAverage;type:int;not null;" json:"dishRatingTagsAverage" `
	CafeteriaID           int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"cafeteriaID"`
	TagID                 int64   `gorm:"column:tagID;foreignKey:tagID;type:int;not null;" json:"tagID"`
	DishID                int64   `gorm:"column:dishID;foreignKey:dishID;type:int;not null;" json:"dishID"`
	Average               float32 `gorm:"column:average;type:float;not null;" json:"average" `
	Min                   int8    `gorm:"column:min;type:int;not null;" json:"min"`
	Max                   int8    `gorm:"column:max;type:int;not null;" json:"max"`
	Std                   float32 `gorm:"column:std;type:float;not null;" json:"std"`
}

// TableName sets the insert table name for this struct type
func (n *InitialDishRatingTagAverage) TableName() string {
	return "dish_rating_tag_average"
}

// migrate20220713000000
// - adds all canteen related models
func migrate20220713000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20220713000000",
		Migrate: func(tx *gorm.DB) error {

			if err := tx.AutoMigrate(
				&InitialCafeteria{},
				&InitialCafeteriaRating{},
				&InitialCafeteriaRatingTag{},
				&InitialCafeteriaRatingTagOption{},
				&InitialDish{},
				&InitialDishesOfTheWeek{},
				&InitialDishNameTagOption{},
				&InitialDishNameTagOptionIncluded{},
				&InitialDishNameTagOptionExcluded{},
				&InitialDishNameTag{},
				&InitialDishRating{},
				&InitialDishRatingTag{},
				&InitialDishRatingTagOption{},
				&InitialDishToDishNameTag{},
				&InitialCafeteriaRatingAverage{},
				&InitialCafeteriaRatingTagAverage{},
				&InitialDishNameTagAverage{},
				&InitialDishRatingAverage{},
				&InitialDishRatingTagAverage{},
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
