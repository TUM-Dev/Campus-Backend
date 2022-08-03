package model

// DishRatingTagOption stores all available options for tags which can be used to quickly rate dishes
type DishRatingTagOption struct {
	DishRatingTagOption int32  `gorm:"primary_key;AUTO_INCREMENT;column:dishRatingTagOption;type:int;not null;" json:"dishRatingTagOption"`
	DE                  string `gorm:"column:DE;type:text;default:de;not null;" json:"DE"`
	EN                  string `gorm:"column:EN;type:text;default:en;not null;" json:"EN"`
}

// TableName sets the insert table name for this struct type
func (n *DishRatingTagOption) TableName() string {
	return "dish_rating_tag_option"
}
