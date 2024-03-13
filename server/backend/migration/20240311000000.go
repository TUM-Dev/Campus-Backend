package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrate20240311000000
// made sure that dishes have the correct indexes
// changed how `dish_rating` is bound to `dish`
func migrate20240311000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240311000000",
		Migrate: func(tx *gorm.DB) error {
			// make sure that dish_rating is FK-bound to dish
			if err := tx.Exec(`alter table dish_rating
					add constraint dish_rating_dish_dish_fk
					foreign key (dishID) references dish (dish)
					on update cascade
				    on delete cascade;`).Error; err != nil {
				return err
			}
			// because dishes already have a cafeteria, storing this again is not necessary
			if err := tx.Exec("alter table dish_rating drop column cafeteriaID").Error; err != nil {
				return err
			}
			// uniqueness
			return tx.Exec("create unique index dish_name_cafeteriaID_uindex on dish (name, cafeteriaID)").Error
		},
		Rollback: func(tx *gorm.DB) error {
			// make sure that dish_rating is FK-bound to dishes
			if err := tx.Exec("alter table dish_rating drop constraint dish_rating_dish_dish_fk").Error; err != nil {
				return err
			}
			// because dishes already have a cafeteria, storing this agiain is not nessesary
			if err := tx.Exec("alter table dish_rating add column cafeteriaID int not null").Error; err != nil {
				return err
			}
			// uniqueness
			return tx.Exec("drop index dish_name_cafeteriaID_uindex on dish").Error
		},
	}
}
