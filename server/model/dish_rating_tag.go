package model

type DishRatingTag struct {
	DishRatingTag int64 `gorm:"primary_key;AUTO_INCREMENT;column:dishRatingTag"`
	RatingID      int64 `gorm:"column:parentRating;not null"`
	//Rating        CanteenRating
	TagID int64 `gorm:"column:tagID;not null"`
	//Tag           CanteenRatingTagOption
	Points int32 `gorm:"not null"`
}
