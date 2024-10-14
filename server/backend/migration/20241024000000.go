package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type newStudentCouncil struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100)"`
	Language    string `gorm:"type:enum('German','English');default:'German'"`
	Description string
}

// TableName sets the insert table name for this struct type
func (n *newStudentCouncil) TableName() string {
	return "student_councils"
}

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

// migrate20241024000000
// - made sure that student councils are supported
func migrate20241024000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20241024000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Migrator().AutoMigrate(&newStudentCouncilCollection{}, &newStudentCouncil{}); err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable(&newStudentCouncil{}, &newStudentCouncilCollection{}); err != nil {
				return err
			}
			return nil
		},
	}
}
