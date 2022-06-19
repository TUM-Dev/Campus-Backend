package cafeteria_rating_models

type MealNameTagOptions struct {
	Id     int32  `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id"`
	NameDE string `gorm:"column:nameDE;type:varchar;size:32" json:"nameDE"`
	NameEN string `gorm:"column:nameEN;type:varchar;size:32" json:"nameEN"`
}

// TableName sets the insert table name for this struct type
func (n *MealNameTagOptions) TableName() string {
	return "meal_name_tag_options"
}
