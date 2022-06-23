package cafeteria_rating_models

// TagRating struct is a row record of the either the meal_tag_rating-table or the cafeteria_rating_tags-table in the database
type CafeteriaRatingTags struct {
	Id           int32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id" `
	ParentRating int32 `gorm:"foreignKey:cafeteriaRatingID;column:parentRating;type:int;" json:"parentRating"`
	Rating       int32 `gorm:"column:rating;type:int;" json:"rating"`
	TagID        int   `gorm:"foreignKey:tagRatingID;column:tagID;type:int" json:"tagID"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRatingTags) TableName() string {
	return "cafeteria_rating_tags"
}
