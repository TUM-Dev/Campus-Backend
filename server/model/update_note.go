package model

// UpdateNote struct for update notes
type UpdateNote struct {
	VersionCode int64  `gorm:"primary_key;autoIncrement;column:version_code;type:int;"`
	VersionName string `gorm:"column:version_name;type:text;"`
	Message     string `gorm:"column:message;type:text;"`
}
