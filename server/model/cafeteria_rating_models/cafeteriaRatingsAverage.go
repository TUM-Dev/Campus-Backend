package cafeteria_rating_models

// News struct is a row record of the mensa table in the tca database
type CafeteriaRatingsAverage struct {
	Id          int32   `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id" `
	CafeteriaID int32   `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	Average     float32 `gorm:"column:average;type:float;" json:"average" `
	Min         int8    `gorm:"column:min;type:int;" json:"min"`
	Max         int8    `gorm:"column:max;type:int;" json:"max"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRatingsAverage) TableName() string {
	return "cafeteria_rating_results"
}
