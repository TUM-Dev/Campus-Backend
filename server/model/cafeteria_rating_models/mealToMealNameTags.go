package cafeteria_rating_models

type MealToMealNameTags struct {
	Id        int32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id"`
	MealID    int32 `gorm:"column:mealID;foreignKey:id;type:int;" json:"mealID"`
	NameTagID int32 `gorm:"foreignKey:Id;column:nameTagID;type:int" json:"nameTagID"`
}

// TableName sets the insert table name for this struct type
func (n *MealToMealNameTags) TableName() string {
	return "meal_to_meal_name_tags"
}
