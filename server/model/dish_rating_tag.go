package model

type DishRatingTag struct {
	DishRatingTag int64 `gorm:"primary_key;autoIncrement;column:dishRatingTag"`
	RatingID      int64 `gorm:"column:parentRating;not null;index:dish_rating_tag_dish_rating_dishRating_fk"`
	//Rating        CanteenRating
	TagID int64 `gorm:"column:tagID;not null;index:dish_rating_tag_dish_rating_tag_option_dishRatingTagOption_fk"`
	//Tag           CanteenRatingTagOption
	Points int32 `gorm:"not null"`
}
