package cafeteria_rating_models

type MealNameTags struct {
	Id           int32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id"`
	ParentRating int32 `gorm:"foreignKey:maelTagID;column:parentRating;type:int;" json:"parentRating"`
	Rating       int32 `gorm:"column:rating;type:int;" json:"rating"`
	TagID        int   `gorm:"foreignKey:tagRatingID;column:tagID;type:int" json:"tagID"`
}

// TableName sets the insert table name for this struct type
func (n *MealNameTags) TableName() string {
	return "meal_name_tags"
}
