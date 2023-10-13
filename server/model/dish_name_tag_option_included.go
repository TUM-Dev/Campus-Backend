package model

type DishNameTagOptionIncluded struct {
	DishNameTagOptionIncluded int64  `gorm:"primary_key;AUTO_INCREMENT;column:dishNameTagOptionIncluded;type:int;not null;" json:"dishNameTagOptionIncluded"`
	NameTagID                 int64  `gorm:"foreignKey:dishNameTagOption;column:nameTagID;type:int;not null;" json:"nameTagID"`
	Expression                string `gorm:"column:expression;type:text;" json:"expression"`
}

// TableName sets the insert table name for this struct type
func (n *DishNameTagOptionIncluded) TableName() string {
	return "dish_name_tag_option_included"
}
