package model

type DishToDishNameTag struct {
	DishToDishNameTag int64 `gorm:"primary_key;AUTO_INCREMENT;column:dishToDishNameTag;type:int;not null;" json:"dishToDishNameTag"`
	DishID            int64 `gorm:"column:dishID;foreignKey:dish;type:int;not null;" json:"dishID"`
	NameTagID         int64 `gorm:"foreignKey:dishNameTagOption;column:nameTagID;type:int;not null;" json:"nameTagID"`
}
