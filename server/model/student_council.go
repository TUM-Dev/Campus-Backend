package model

import (
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// StudentCouncil stores a student Council
type StudentCouncil struct {
	gorm.Model
	Name                       string
	Language                   string `gorm:"type:enum('German','English');default:'German'"`
	Description                null.String
	LinkUrl                    null.String `gorm:"type:varchar(190);unique;"`
	ImageID                    null.Int
	Image                      *File `gorm:"foreignKey:ImageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ImageCaption               null.String
	StudentCouncilCollectionID uint
	StudentCouncilCollection   StudentCouncilCollection `gorm:"foreignKey:StudentCouncilCollectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
