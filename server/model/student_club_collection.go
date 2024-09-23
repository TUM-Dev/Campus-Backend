package model

import "gorm.io/gorm"

// StudentClubCollection stores what collection a club belongs to
type StudentClubCollection struct {
	gorm.Model
	Name        string `gorm:"primaryKey;type:varchar(100)"`
	Language    string `gorm:"type:enum('German','English');default:'German'"`
	Description string
}
