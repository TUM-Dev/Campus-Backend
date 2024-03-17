package migration

import (
	"fmt"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type missingFK struct {
	fromTable  string
	fromColumn string
	toTable    string
	toColumn   string
}

func missingDishFKs() []missingFK {
	return []missingFK{
		{"cafeteria_rating", "cafeteriaID", "cafeteria", "cafeteria"},
		{"dish", "cafeteriaID", "cafeteria", "cafeteria"},
		{"dish2dishflags", "dish", "dish", "dish"},
		{"dish2dishflags", "flag", "dishflags", "flag"},
		{"dish2mensa", "dish", "dish", "dish"},
		{"dish2mensa", "mensa", "mensa", "mensa"},
		{"dish_rating_tag", "parentRating", "dish_rating", "dishRating"},
		{"dish_rating_tag", "tagID", "dish_rating_tag_option", "dishRatingTagOption"},
		{"dish_to_dish_name_tag", "dishID", "dish", "dish"},
		{"dishes_of_the_week", "dishID", "dish", "dish"},
	}
}

// migrate20240320000000
// made sure that all dish ratings are connected with FK-Relationships
func migrate20240320000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240320000000",
		Migrate: func(tx *gorm.DB) error {
			for _, fk := range missingDishFKs() {
				tx.Exec(fmt.Sprintf("alter table `%s` add constraint %s_%s_%s_fk foreign key (`%s`) references `%s` (`%s`) on update cascade on delete cascade",
					fk.fromTable, fk.fromTable, fk.toTable, fk.toColumn, fk.fromColumn, fk.toTable, fk.toColumn))
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			for _, fk := range missingDishFKs() {
				tx.Exec(fmt.Sprintf("alter table `%s` drop foreign key %s_%s_%s_fk",
					fk.fromTable, fk.fromTable, fk.toTable, fk.toColumn))
			}
			return nil
		},
	}
}
