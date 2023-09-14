package model

// UpdateNote struct for update notes
type UpdateNote struct {
	VersionCode int32  `gorm:"primary_key;AUTO_INCREMENT;column:version_code;type:int;"`
	VersionName string `gorm:"column:version_name;type:text;"`
	Message     string `gorm:"column:message;type:text;"`
}

// TableName sets the insert table name for this struct type
func (n *UpdateNote) TableName() string {
	return "update_note"
}
