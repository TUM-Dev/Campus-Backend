package cafeteria_rating_models

type MealNameTagOption struct {
	MealNameTagOption int32  `gorm:"primary_key;AUTO_INCREMENT;column:mealNameTagOption;type:int;" json:"mealNameTagOption"`
	DE                string `gorm:"column:DE;type:text" json:"DE"`
	EN                string `gorm:"column:EN;type:text" json:"EN"`
}

// TableName sets the insert table name for this struct type
func (n *MealNameTagOption) TableName() string {
	return "meal_name_tag_option"
}
