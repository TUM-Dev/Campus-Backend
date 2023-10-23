package model

type IncludedDishNameTagOption struct {
	DishNameTagOptionIncluded int64  `gorm:"primary_key;AUTO_INCREMENT;column:dishNameTagOptionIncluded;type:int;not null;" json:"dishNameTagOptionIncluded"`
	NameTagID                 int64  `gorm:"foreignKey:dishNameTagOption;column:nameTagID;type:int;not null;" json:"nameTagID"`
	Expression                string `gorm:"column:expression;type:text;" json:"expression"`
}
