package model

type DishesOfTheWeek struct {
	DishesOfTheWeek int64 `gorm:"primary_key;AUTO_INCREMENT;column:dishesOfTheWeek;type:int;not null;" json:"dishesOfTheWeek"`
	Year            int32 `gorm:"column:year;type:int;not null;" json:"year"`
	Week            int32 `gorm:"column:week;type:int;not null;" json:"week"`
	Day             int32 `gorm:"column:day;type:int;not null;" json:"day"`
	DishID          int64 `gorm:"column:dishID;foreignKey:dish;type:int;not null;" json:"dishID"`
}

// TableName sets the insert table name for this struct type
func (n *DishesOfTheWeek) TableName() string {
	return "dishes_of_the_week"
}
