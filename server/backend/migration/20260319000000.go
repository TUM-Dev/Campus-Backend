package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20260319000000
// - changes the unique constraint on student_clubs.link_url to a composite
//   unique constraint on (language, link_url) so the same club URL can exist
//   for both the German and English language variants
func migrate20260319000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260319000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec("ALTER TABLE `student_clubs` DROP INDEX `uni_student_clubs_link_url`").Error; err != nil {
				return err
			}
			return tx.Exec("ALTER TABLE `student_clubs` ADD UNIQUE INDEX `uni_student_clubs_language_link_url` (`language`, `link_url`)").Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec("ALTER TABLE `student_clubs` DROP INDEX `uni_student_clubs_language_link_url`").Error; err != nil {
				return err
			}
			return tx.Exec("ALTER TABLE `student_clubs` ADD UNIQUE INDEX `uni_student_clubs_link_url` (`link_url`)").Error
		},
	}
}
