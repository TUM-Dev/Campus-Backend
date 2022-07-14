package cafeteria_rating_models

// MealRatingTagOption stores all available options for tags which can be used to quickly rate meals
type MealRatingTagOption struct {
	MealRatingTagOption int32  `gorm:"primary_key;AUTO_INCREMENT;column:mealRatingTagOption;type:int;" json:"mealRatingTagOption"`
	DE                  string `gorm:"column:DE;type:mediumtext;default:de" json:"DE"`
	EN                  string `gorm:"column:EN;type:mediumtext;default:en" json:"EN"`
}

// TableName sets the insert table name for this struct type
func (n *MealRatingTagOption) TableName() string {
	return "meal_rating_tag_option"
}
