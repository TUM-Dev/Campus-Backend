package cafeteria_rating_models

type DishRatingTag struct {
	DishRatingTag       int32 `gorm:"primary_key;AUTO_INCREMENT;column:dishRatingTag;type:int;" json:"dishRatingTag"`
	CorrespondingRating int32 `gorm:"foreignKey:cafeteriaRating;column:parentRating;type:int;" json:"parentRating"`
	Points              int32 `gorm:"column:points;type:int;" json:"points"`
	TagID               int   `gorm:"foreignKey:dishRatingTagOption;column:tagID;type:int" json:"tagID"`
}

// TableName sets the insert table name for this struct type
func (n *DishRatingTag) TableName() string {
	return "dish_rating_tag"
}
