package migration

import (
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func kinoTablesWithWrongNullability() []string {
	return []string{
		"year varchar(4)",
		"runtime varchar(40)",
		"genre varchar(100)",
		"director text",
		"actors text",
		"rating varchar(4)",
	}
}

// migrate20240402000000
// Aligned a few nullability issues between the data scraped from the tufilm and our db
func migrate20240402000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240402000000",
		Migrate: func(tx *gorm.DB) error {
			for _, f := range kinoTablesWithWrongNullability() {
				if err := tx.Exec(fmt.Sprintf("alter table kino modify %s null", f)).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			for _, f := range kinoTablesWithWrongNullability() {
				if err := tx.Exec(fmt.Sprintf("alter table kino modify %s not null", f)).Error; err != nil {
					return err
				}
			}
			return nil
		},
	}
}
