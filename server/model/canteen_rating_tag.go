package model

// CanteenRatingTag struct is a row record of the either the dish_tag_ratings-table or the cafeteria_rating_tags-table in the database
type CanteenRatingTag struct {
	CafeteriaRatingTag  int64 `gorm:"primary_key;AUTO_INCREMENT;column:CanteenRatingTag;type:int;not null;" json:"CanteenRatingTag" `
	CorrespondingRating int64 `gorm:"foreignKey:cafeteriaRatingID;column:correspondingRating;type:int;not null;" json:"correspondingRating"`
	Points              int32 `gorm:"column:points;type:int;not null;" json:"points"`
	TagID               int64 `gorm:"foreignKey:cafeteriaRatingTagOption;column:tagID;type:int;not null;" json:"tagID"`
}
