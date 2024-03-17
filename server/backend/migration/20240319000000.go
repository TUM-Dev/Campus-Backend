package migration

import (
	"fmt"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20240319000000
// made sure that all timestamp columns are defaulted and update as expected
func migrate20240319000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240319000000",
		Migrate: func(tx *gorm.DB) error {
			tables := []string{"cafeteria_rating", "canteen_head_count", "dish_rating"}
			for _, t := range tables {
				tx.Exec(fmt.Sprintf("alter table `%s` modify timestamp timestamp default current_timestamp() not null on update current_timestamp()", t))
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			tables := []string{"cafeteria_rating", "canteen_head_count", "dish_rating"}
			for _, t := range tables {
				tx.Exec(fmt.Sprintf("alter table `%s` modify timestamp timestamp not null", t))
			}
			return nil
		},
	}
}
