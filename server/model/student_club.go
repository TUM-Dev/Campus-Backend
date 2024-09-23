package model

import (
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// StudentClub stores a student Club
type StudentClub struct {
	gorm.Model
	Name                    string
	Language                string `gorm:"type:enum('German','English');default:'German'"`
	Description             null.String
	LinkUrl                 null.String `gorm:"type:varchar(190);unique;"`
	ImageID                 null.Int
	Image                   *File `gorm:"foreignKey:ImageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ImageCaption            null.String
	StudentClubCollectionID string                `gorm:"type:varchar(100)"`
	StudentClubCollection   StudentClubCollection `gorm:"foreignKey:StudentClubCollectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
