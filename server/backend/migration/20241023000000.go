package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20241023000000
// - made sure that student clubs are localised
func migrate20241023000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20241023000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec("alter table student_clubs add language enum ('German', 'English') default 'German' not null").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table student_club_collections add language enum ('German', 'English') default 'German' not null").Error; err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec("alter table student_clubs drop column language").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table student_club_collections drop column language").Error; err != nil {
				return err
			}
			return nil
		},
	}
}
