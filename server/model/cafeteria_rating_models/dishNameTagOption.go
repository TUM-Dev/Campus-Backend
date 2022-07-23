package cafeteria_rating_models

type DishNameTagOption struct {
	DishNameTagOption int32  `gorm:"column:dishNameTagOption;type:int;primary_key;AUTO_INCREMENT;" json:"dishNameTagOption"`
	DE                string `gorm:"column:DE;type:text" json:"DE"`
	EN                string `gorm:"column:EN;type:text" json:"EN"`
}

// TableName sets the insert table name for this struct type
func (n *DishNameTagOption) TableName() string {
	return "dish_name_tag_option"
}
