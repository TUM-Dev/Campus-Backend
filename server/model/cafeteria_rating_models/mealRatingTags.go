package cafeteria_rating_models

type MealTagRating struct {
	Id           int32  `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id" :"id"`
	ParentRating int32  `gorm:"foreignKey:cafeteriaRatingID;column:parentRating;type:int;" json:"parentRating" :"id"`
	Rating       int32  `gorm:"column:rating;type:int;" json:"rating" :"rating"`
	TagID        string `gorm:"foreignKey:cafeteriaRatingID;column:tagID;type:int" json:"tagID" :"tagname"`
}

// TableName sets the insert table name for this struct type
func (n *MealTagRating) TableName() string {
	return "meal_rating_tags_options"
}
