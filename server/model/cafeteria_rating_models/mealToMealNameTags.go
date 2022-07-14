package cafeteria_rating_models

type MealToMealNameTag struct {
	MealToMealNameTag int32 `gorm:"primary_key;AUTO_INCREMENT;column:mealToMealNameTag;type:int;" json:"mealToMealNameTag"`
	MealID            int32 `gorm:"column:mealID;foreignKey:meal;type:int;" json:"mealID"`
	NameTagID         int32 `gorm:"foreignKey:mealNameTagOption;column:nameTagID;type:int" json:"nameTagID"`
}

// TableName sets the insert table name for this struct type
func (n *MealToMealNameTag) TableName() string {
	return "meal_to_meal_name_tag"
}
