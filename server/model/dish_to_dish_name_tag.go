package model

type DishToDishNameTag struct {
	DishToDishNameTag int64 `gorm:"primary_key;autoIncrement;column:dishToDishNameTag;type:int;not null;" json:"dishToDishNameTag"`
	DishID            int64 `gorm:"column:dishID;foreignKey:dish;type:int;not null;index:dish_to_dish_name_tag_dish_dish_fk;" json:"dishID"`
	NameTagID         int64 `gorm:"foreignKey:dishNameTagOption;column:nameTagID;type:int;not null;" json:"nameTagID"`
}
