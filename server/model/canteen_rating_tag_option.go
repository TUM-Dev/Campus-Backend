package model

// CanteenRatingTagOption stores all available options for tags which can be used to quickly rate cafeterias
type CanteenRatingTagOption struct {
	CafeteriaRatingsTagOption int64  `gorm:"primary_key;autoIncrement;column:cafeteriaRatingTagOption;type:int;not null;" json:"canteenRatingTagOption"`
	DE                        string `gorm:"column:DE;text;default:('de');not null;" json:"DE"`
	EN                        string `gorm:"column:EN;text;default:('en');not null;" json:"EN"`
}
