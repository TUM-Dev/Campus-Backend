package model

// CafeteriaRatingTag struct is a row record of the either the dish_tag_rating-table or the cafeteria_rating_tags-table in the database
type CafeteriaRatingTag struct {
	CafeteriaRatingTag  int32 `gorm:"primary_key;AUTO_INCREMENT;column:CafeteriaRatingTag;type:int;not null;" json:"CanteenRatingTag" `
	CorrespondingRating int32 `gorm:"foreignKey:cafeteriaRatingID;column:correspondingRating;type:int;not null;" json:"correspondingRating"`
	Points              int32 `gorm:"column:points;type:int;not null;" json:"points"`
	TagID               int   `gorm:"foreignKey:cafeteriaRatingTagOption;column:tagID;type:int;not null;" json:"tagID"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRatingTag) TableName() string {
	return "cafeteria_rating_tag"
}
