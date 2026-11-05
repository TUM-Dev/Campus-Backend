package model

import (
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// StudentCouncil stores a student council (e.g. a Fachschaft)
type StudentCouncil struct {
	gorm.Model
	Name                       string
	Language                   string `gorm:"type:enum('German','English');default:'German';uniqueIndex:uni_student_councils_link_url"`
	Description                null.String
	LinkUrl                    null.String `gorm:"type:varchar(190);uniqueIndex:uni_student_councils_link_url"`
	ImageID                    null.Int
	Image                      *File `gorm:"foreignKey:ImageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ImageCaption               null.String
	StudentCouncilCollectionID uint
	StudentCouncilCollection   StudentCouncilCollection `gorm:"foreignKey:StudentCouncilCollectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
