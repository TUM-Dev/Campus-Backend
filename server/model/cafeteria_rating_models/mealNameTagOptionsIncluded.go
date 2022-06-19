package cafeteria_rating_models

type MealNameTagOptionsIncluded struct {
	Id         int32  `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id"`
	NameTagID  int32  `gorm:"foreignKey:Id;column:nameTagID;type:int" json:"nameTagID"`
	Expression string `gorm:"column:expression;type:mediumtext" json:"expression"`
}

// TableName sets the insert table name for this struct type
func (n *MealNameTagOptionsIncluded) TableName() string {
	return "meal_name_tag_options_included"
}
