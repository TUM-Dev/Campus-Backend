package model

import "gorm.io/gorm"

// StudentCouncilCollection stores which collection a council belongs to
type StudentCouncilCollection struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100)"`
	Language    string `gorm:"type:enum('German','English');default:'German'"`
	Description string
}
