package model

// CanteenRatingTag connects Canteens with ratings
type CanteenRatingTag struct {
	CanteenRatingTag      int64 `gorm:"primary_key;AUTO_INCREMENT;column:CafeteriaRatingTag"`
	CorrespondingRatingID int64
	CorrespondingRating   CanteenRating
	TagID                 int64
	Tag                   CanteenRatingTagOption
	Points                int32
}
