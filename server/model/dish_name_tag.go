package model

type DishNameTag struct {
	DishNameTag int64 `gorm:"primary_key;AUTO_INCREMENT;column:DishNameTag;not null;"`
	RatingID    int64 `gorm:"foreignKey:dish;column:correspondingRating;not null;"`
	//Rating      CanteenRating
	TagNameID int64 `gorm:"foreignKey:tagRatingID;column:tagNameID;not null;"`
	//TagName     DishNameTag
	Points int32 `gorm:"not null;"`
}
