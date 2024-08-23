package model

// CanteenRatingTag connects Canteens with ratings
type CanteenRatingTag struct {
	CanteenRatingTag int64 `gorm:"primary_key;AUTO_INCREMENT;column:CafeteriaRatingTag"`
	RatingID         int64 `gorm:"column:corresponding_rating_id"`
	//Rating           CanteenRating
	TagID int64 `gorm:"column:tag_id"`
	//Tag              CanteenRatingTagOption
	Points int32 `gorm:"not null"`
}
