package cafeteria_rating_models

type MealRatingTag struct {
	MealRatingTag       int32 `gorm:"primary_key;AUTO_INCREMENT;column:mealRatingTag;type:int;" json:"mealRatingTag"`
	CorrespondingRating int32 `gorm:"foreignKey:cafeteriaRatingID;column:parentRating;type:int;" json:"parentRating"`
	Points              int32 `gorm:"column:points;type:int;" json:"points"`
	TagID               int   `gorm:"foreignKey:tagRatingID;column:tagID;type:int" json:"tagID"`
}

// TableName sets the insert table name for this struct type
func (n *MealRatingTag) TableName() string {
	return "meal_rating_tag"
}
