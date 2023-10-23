package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type wrongTableName struct {
	Original string
	New      string
}

var wrongTableNames = []wrongTableName{
	{"crontab", "crontabs"},
	{"kino", "movies"},
	{"dish_name_tag_option_included", "included_dish_name_tag_options"},
	{"dish_name_tag_option", "dish_name_tag_options"},
	{"cafeteria", "canteens"},
	{"cafeteria_rating", "canteen_ratings"},
	{"cafeteria_rating_tag", "canteen_rating_tags"},
	{"cafeteria_rating_tag_option", "canteen_rating_tag_options"},
	{"dish", "dishes"},
	{"dish_name_tag", "dish_name_tags"},
	{"dish_name_tag_option_excluded", "excluded_dish_name_tag_options"},
	{"dish_name_tag_option_included", "included_dish_name_tag_option"},
}

// migrate20231024000000
// - replaces all instances of misleadingly named tables with the correct ones
func migrate20231024000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20231024000000",
		Migrate: func(tx *gorm.DB) error {
			for _, table := range wrongTableNames {
				if err := tx.Migrator().RenameTable(table.Original, table.New); err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			for _, table := range wrongTableNames {
				if err := tx.Migrator().RenameTable(table.New, table.Original); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
