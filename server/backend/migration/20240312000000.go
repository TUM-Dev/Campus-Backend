package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20240212000000
// implemented a basic variant of spam protection
func migrate20240212000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240212000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec("alter table feedback modify email_id text charset utf8mb3 not null").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table feedback modify receiver text charset utf8mb3 not null").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table feedback modify feedback text charset utf8mb3 not null").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table feedback modify image_count int not null").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table feedback modify timestamp datetime default current_timestamp() not null").Error; err != nil {
				return err
			}
			return tx.Exec("create unique index receiver_reply_to_feedback_app_version_uindex on feedback (receiver,reply_to,feedback,app_version)").Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec("alter table feedback modify email_id text charset utf8mb3 null").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table feedback modify receiver text charset utf8mb3 null").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table feedback modify feedback text charset utf8mb3 null").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table feedback modify image_count int null").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table feedback modify timestamp datetime default current_timestamp() null").Error; err != nil {
				return err
			}
			return tx.Exec("drop index receiver_reply_to_feedback_app_version_uindex on feedback").Error
		},
	}
}
