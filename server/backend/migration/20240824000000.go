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
	//{"crontab", "crontabs"}, <- too much work for now
	{"kino", "movies"},
	{"cafeteria", "canteens"},
	{"canteen_head_count", "canteen_head_counts"},
	{"cafeteria_rating", "canteen_ratings"},
	{"cafeteria_rating_tag", "canteen_rating_tags"},
	{"cafeteria_rating_tag_option", "canteen_rating_tag_options"},
	{"dish", "dishes"},
	{"dish_rating", "dish_ratings"},
	{"dish_rating_tag", "dish_rating_tags"},
	{"dish_rating_tag_option", "dish_rating_tag_options"},
	{"dish_name_tag", "dish_name_tags"},
	{"dish_name_tag_option", "dish_name_tag_options"},
	{"dish_name_tag_option_excluded", "excluded_dish_name_tag_options"},
	{"dish_name_tag_option_included", "included_dish_name_tag_options"},
	{"dish_to_dish_name_tag", "dish_to_dish_name_tags"},
	{"dishes_of_the_week", "dishes_of_the_weeks"},
	{"update_note", "update_notes"},
	{"newsSource", "news_sources"},
	{"news_alert", "news_alerts"},
}

// migrate20240824000000
// - replaces all instances of misleadingly named tables with the correct ones
func migrate20240824000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240824000000",
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
