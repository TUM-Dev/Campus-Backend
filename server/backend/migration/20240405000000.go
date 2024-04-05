package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20240405000000
// Split the reply-to fields of the feedback model into name and email
func migrate20240405000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240405000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec("alter table feedback change reply_to reply_to_email text null").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table feedback add reply_to_name text null default null after reply_to_email").Error; err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec("alter table feedback change reply_to_email reply_to text null").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table feedback drop column reply_to_name").Error; err != nil {
				return err
			}
			return nil
		},
	}
}
