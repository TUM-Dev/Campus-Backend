package cafeteria_rating_models

// TagRating struct is a row record of the either the meal_tag_rating-table or the cafeteria_rating_tags-table in the database
type CafeteriaRatingTag struct {
	CafeteriaRatingTag  int32 `gorm:"primary_key;AUTO_INCREMENT;column:CafeteriaRatingTag;type:int;" json:"CafeteriaRatingTag" `
	CorrespondingRating int32 `gorm:"foreignKey:cafeteriaRatingID;column:correspondingRating;type:int;" json:"correspondingRating"`
	Points              int32 `gorm:"column:points;type:int;" json:"points"`
	TagID               int   `gorm:"foreignKey:tagRatingID;column:tagID;type:int" json:"tagID"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRatingTag) TableName() string {
	return "cafeteria_rating_tag"
}
