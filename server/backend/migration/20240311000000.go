package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20240311000000
// made sure that dishes have the correct indexes
func migrate20240311000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240311000000",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec("create unique index dish_name_cafeteriaID_uindex on dish (name, cafeteriaID)").Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec("drop index dish_name_cafeteriaID_uindex on dish").Error
		},
	}
}
