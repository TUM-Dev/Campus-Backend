package cafeteria_rating_models

type DishToDishNameTag struct {
	DishToDishNameTag int32 `gorm:"primary_key;AUTO_INCREMENT;column:dishToDishNameTag;type:int;" json:"dishToDishNameTag"`
	DishID            int32 `gorm:"column:dishID;foreignKey:dish;type:int;" json:"dishID"`
	NameTagID         int32 `gorm:"foreignKey:dishNameTagOption;column:nameTagID;type:int" json:"nameTagID"`
}

// TableName sets the insert table name for this struct type
func (n *DishToDishNameTag) TableName() string {
	return "dish_to_dish_name_tag"
}
