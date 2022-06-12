package cafeteria_rating_models

type MealRatingsTags struct {
	Id           int32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id"`
	ParentRating int32 `gorm:"foreignKey:cafeteriaRatingID;column:parentRating;type:int;" json:"parentRating"`
	Rating       int32 `gorm:"column:rating;type:int;" json:"rating"`
	TagID        int   `gorm:"foreignKey:tagRatingID;column:tagID;type:int" json:"tagID"`
}

// TableName sets the insert table name for this struct type
func (n *MealRatingsTags) TableName() string {
	return "meal_rating_tags_options"
}
