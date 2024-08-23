package model

type DishRatingTag struct {
	DishRatingTag int64 `gorm:"primary_key;AUTO_INCREMENT;column:dishRatingTag"`
	RatingID      int64 `gorm:"column:parentRating"`
	Rating        CanteenRating
	TagID         int64
	Tag           CanteenRatingTagOption
	Points        int32 `gorm:"not null"`
}
