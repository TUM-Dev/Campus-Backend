package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

type newStudentCouncilCollection struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100)"`
	Language    string `gorm:"type:enum('German','English');default:'German'"`
	Description string
}

// TableName sets the insert table name for this struct type
func (n *newStudentCouncilCollection) TableName() string {
	return "student_council_collections"
}

type newStudentCouncil struct {
	gorm.Model
	Name                       string
	Language                   string `gorm:"type:enum('German','English');default:'German';uniqueIndex:uni_student_councils_link_url"`
	Description                null.String
	LinkUrl                    null.String `gorm:"type:varchar(190);uniqueIndex:uni_student_councils_link_url"`
	ImageID                    null.Int
	Image                      *File `gorm:"foreignKey:ImageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ImageCaption               null.String
	StudentCouncilCollectionID uint
	StudentCouncilCollection   newStudentCouncilCollection `gorm:"foreignKey:StudentCouncilCollectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName sets the insert table name for this struct type
func (n *newStudentCouncil) TableName() string {
	return "student_councils"
}

// migrate20261105000000
// - adds the student_councils and student_council_collections tables
func migrate20261105000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20261105000000",
		Migrate: func(tx *gorm.DB) error {
			return tx.Migrator().AutoMigrate(&newStudentCouncilCollection{}, &newStudentCouncil{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&newStudentCouncil{}, &newStudentCouncilCollection{})
		},
	}
}
