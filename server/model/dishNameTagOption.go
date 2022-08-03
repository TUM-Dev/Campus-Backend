package model

type DishNameTagOption struct {
	DishNameTagOption int32  `gorm:"column:dishNameTagOption;type:int;primary_key;AUTO_INCREMENT;not null;" json:"dishNameTagOption"`
	DE                string `gorm:"column:DE;type:text;not null;" json:"DE"`
	EN                string `gorm:"column:EN;type:text;not null;" json:"EN"`
}

// TableName sets the insert table name for this struct type
func (n *DishNameTagOption) TableName() string {
	return "dish_name_tag_option"
}
