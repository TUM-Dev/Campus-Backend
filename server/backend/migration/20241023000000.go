package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// StudentClub stores a student Club
type newStudentClub struct {
	gorm.Model
	Name                    string
	Language                string `gorm:"type:enum('German','English');default:'German'"`
	Description             null.String
	LinkUrl                 null.String `gorm:"type:varchar(190);unique;"`
	ImageID                 null.Int
	Image                   *File `gorm:"foreignKey:ImageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ImageCaption            null.String
	StudentClubCollectionID uint
	StudentClubCollection   newStudentClubCollection `gorm:"foreignKey:StudentClubCollectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName sets the insert table name for this struct type
func (n *newStudentClub) TableName() string {
	return "student_clubs"
}

type newStudentClubCollection struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100)"`
	Language    string `gorm:"type:enum('German','English');default:'German'"`
	Description string
}

// TableName sets the insert table name for this struct type
func (n *newStudentClubCollection) TableName() string {
	return "student_club_collections"
}

// migrate20241023000000
// - made sure that student clubs are localised
func migrate20241023000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20241023000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&InitialStudentClub{}, &InitialStudentClubCollection{}); err != nil {
				return err
			}
			if err := tx.Migrator().AutoMigrate(&newStudentClubCollection{}, &newStudentClub{}); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&newStudentClub{}, &newStudentClubCollection{}); err != nil {
				return err
			}
			if err := tx.Migrator().AutoMigrate(&InitialStudentClubCollection{}, &InitialStudentClub{}); err != nil {
				return err
			}
			return nil
		},
	}
}
