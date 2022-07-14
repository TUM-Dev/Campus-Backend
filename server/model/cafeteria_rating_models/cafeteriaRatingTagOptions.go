package cafeteria_rating_models

// CafeteriaRatingTagOption stores all available options for tags which can be used to quickly rate cafeterias
type CafeteriaRatingTagOption struct {
	CafeteriaRatingsTagOption int32  `gorm:"primary_key;AUTO_INCREMENT;column:cafeteriaRatingTagOption;type:int;" json:"cafeteriaRatingTagOption"`
	DE                        string `gorm:"column:DE;text;default:de" json:"DE"`
	EN                        string `gorm:"column:EN;text;default:en" json:"EN"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRatingTagOption) TableName() string {
	return "cafeteria_rating_tag_option"
}
