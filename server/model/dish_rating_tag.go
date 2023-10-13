package model

type DishRatingTag struct {
	DishRatingTag       int64 `gorm:"primary_key;AUTO_INCREMENT;column:dishRatingTag;type:int;not null;" json:"dishRatingTag"`
	CorrespondingRating int64 `gorm:"foreignKey:cafeteriaRating;column:parentRating;type:int;not null;" json:"parentRating"`
	Points              int32 `gorm:"column:points;type:int;not null;" json:"points"`
	TagID               int64 `gorm:"foreignKey:dishRatingTagOption;column:tagID;type:int;not null;" json:"tagID"`
}

// TableName sets the insert table name for this struct type
func (n *DishRatingTag) TableName() string {
	return "dish_rating_tag"
}
