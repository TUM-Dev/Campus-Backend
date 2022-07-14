package cafeteria_rating_models

type MealNameTagOptionExcluded struct {
	MealNameTagOptionExcluded int32  `gorm:"primary_key;AUTO_INCREMENT;column:mealNameTagOptionExcluded;type:int;" json:"mealNameTagOptionExcluded"`
	NameTagID                 int32  `gorm:"foreignKey:Id;column:nameTagID;type:int" json:"nameTagID"`
	Expression                string `gorm:"column:expression;type:text" json:"expression"`
}

// TableName sets the insert table name for this struct type
func (n *MealNameTagOptionExcluded) TableName() string {
	return "meal_name_tag_option_excluded"
}
