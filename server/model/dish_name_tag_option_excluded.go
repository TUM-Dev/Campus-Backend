package model

type DishNameTagOptionExcluded struct {
	DishNameTagOptionExcluded int64  `gorm:"primary_key;AUTO_INCREMENT;column:dishNameTagOptionExcluded;type:int;not null;" json:"dishNameTagOptionExcluded"`
	NameTagID                 int64  `gorm:"foreignKey:dishNameTagOption;column:nameTagID;type:int;not null;" json:"nameTagID"`
	Expression                string `gorm:"column:expression;type:text;" json:"expression"`
}

// TableName sets the insert table name for this struct type
func (n *DishNameTagOptionExcluded) TableName() string {
	return "dish_name_tag_option_excluded"
}
