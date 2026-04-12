package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20250412000000
// - changes the unique index on student_clubs.link_url from a single-column
//   index to a composite index on (link_url, language) so that the same URL
//   can appear in both German and English rows
func migrate20250412000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20250412000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec("ALTER TABLE `student_clubs` DROP INDEX `uni_student_clubs_link_url`").Error; err != nil {
				return err
			}
			return tx.Exec("CREATE UNIQUE INDEX `uni_student_clubs_link_url` ON `student_clubs` (`link_url`, `language`)").Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec("ALTER TABLE `student_clubs` DROP INDEX `uni_student_clubs_link_url`").Error; err != nil {
				return err
			}
			return tx.Exec("CREATE UNIQUE INDEX `uni_student_clubs_link_url` ON `student_clubs` (`link_url`)").Error
		},
	}
}
