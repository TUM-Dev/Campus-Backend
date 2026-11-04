package model

type ExcludedDishNameTagOption struct {
	DishNameTagOptionExcluded int64  `gorm:"primary_key;autoIncrement;column:dishNameTagOptionExcluded;type:int;not null;" json:"dishNameTagOptionExcluded"`
	NameTagID                 int64  `gorm:"foreignKey:dishNameTagOption;column:nameTagID;type:int;not null;" json:"nameTagID"`
	Expression                string `gorm:"column:expression;type:text;" json:"expression"`
}
