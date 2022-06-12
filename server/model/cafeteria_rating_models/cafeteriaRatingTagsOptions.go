package cafeteria_rating_models

type CafeteriaRatingsTagsOptions struct {
	Id     int32  `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id"`
	NameDE string `gorm:"column:nameDE;type:varchar;size:32" json:"nameDE"`
	NameEN string `gorm:"column:nameEN;type:varchar;size:32" json:"nameEN"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRatingsTagsOptions) TableName() string {
	return "cafeteria_rating_tags_options"
}
