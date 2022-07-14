package cafeteria_rating_models

type MealNameTag struct {
	MealNameTag         int32 `gorm:"primary_key;AUTO_INCREMENT;column:MealNameTag;type:int;" json:"MealNameTag"`
	CorrespondingRating int32 `gorm:"foreignKey:mealTagID;column:correspondingRating;type:int;" json:"correspondingRating"`
	Points              int32 `gorm:"column:points;type:int;" json:"points"`
	TagNameID           int   `gorm:"foreignKey:tagRatingID;column:tagNameID;type:int" json:"tagnameID"`
}

// TableName sets the insert table name for this struct type
func (n *MealNameTag) TableName() string {
	return "meal_name_tag"
}
