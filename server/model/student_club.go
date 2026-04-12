package model

import (
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// StudentClub stores a student Club
type StudentClub struct {
	gorm.Model
	Name                    string
	Language                string      `gorm:"type:enum('German','English');default:'German';uniqueIndex:uni_student_clubs_link_url"`
	Description             null.String
	LinkUrl                 null.String `gorm:"type:varchar(190);uniqueIndex:uni_student_clubs_link_url"`
	ImageID                 null.Int
	Image                   *File `gorm:"foreignKey:ImageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ImageCaption            null.String
	StudentClubCollectionID uint
	StudentClubCollection   StudentClubCollection `gorm:"foreignKey:StudentClubCollectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
