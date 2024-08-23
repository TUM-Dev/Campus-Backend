package model

type DishNameTagOption struct {
	DishNameTagOption int64  `gorm:"primary_key;autoIncrement;column:dishNameTagOption;type:int;not null;" json:"dishNameTagOption"`
	DE                string `gorm:"column:DE;type:text;not null;" json:"DE"`
	EN                string `gorm:"column:EN;type:text;not null;" json:"EN"`
}
