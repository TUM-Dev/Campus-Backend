package model

type DishNameTag struct {
	DishNameTag         int32 `gorm:"primary_key;AUTO_INCREMENT;column:DishNameTag;type:int;not null;" json:"DishNameTag"`
	CorrespondingRating int32 `gorm:"foreignKey:dish;column:correspondingRating;type:int;not null;" json:"correspondingRating"`
	Points              int32 `gorm:"column:points;type:int;not null;" json:"points"`
	TagNameID           int   `gorm:"foreignKey:tagRatingID;column:tagNameID;type:int;not null;" json:"tagnameID"`
}

// TableName sets the insert table name for this struct type
func (n *DishNameTag) TableName() string {
	return "dish_name_tag"
}
