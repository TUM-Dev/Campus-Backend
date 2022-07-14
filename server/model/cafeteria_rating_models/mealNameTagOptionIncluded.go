package cafeteria_rating_models

type MealNameTagOptionIncluded struct {
	MealNameTagOptionIncluded int32  `gorm:"primary_key;AUTO_INCREMENT;column:mealNameTagOptionIncluded;type:int;" json:"mealNameTagOptionIncluded"`
	NameTagID                 int32  `gorm:"foreignKey:Id;column:nameTagID;type:int" json:"nameTagID"`
	Expression                string `gorm:"column:expression;type:text" json:"expression"`
}

// TableName sets the insert table name for this struct type
func (n *MealNameTagOptionIncluded) TableName() string {
	return "meal_name_tag_option_included"
}
