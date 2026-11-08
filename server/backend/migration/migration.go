// Package migration contains functions related to database changes and executes them
package migration

import (
	"fmt"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Models returns all models whose schema is managed by migrations.
// It is shared by AutoMigrate and the Atlas schema loader (tools/atlasloader).
func Models() []any {
	return []any{
		&model.AdminSession{},
		&model.AdminAudit{},
		&model.Canteen{},
		&model.CanteenRating{},
		&model.CanteenRatingTag{},
		&model.CanteenRatingTagOption{},
		&model.CanteenHeadCount{},
		&model.Crontab{},
		&model.Device{},
		&model.Dish{},
		&model.DishNameTag{},
		&model.DishNameTagOption{},
		&model.ExcludedDishNameTagOption{},
		&model.IncludedDishNameTagOption{},
		&model.DishRating{},
		&model.DishRatingTag{},
		&model.DishRatingTagOption{},
		&model.DishToDishNameTag{},
		&model.DishesOfTheWeek{},
		&model.Feedback{},
		&model.File{},
		&model.Movie{},
		&model.NewExamResultsSubscriber{},
		&model.News{},
		&model.NewsAlert{},
		&model.NewsSource{},
		&model.Notification{},
		&model.NotificationConfirmation{},
		&model.NotificationType{},
		&model.UpdateNote{},
		&model.StudentClub{},
		&model.StudentClubCollection{},
		&model.StudentCouncil{},
		&model.StudentCouncilCollection{},
	}
}

// legacyMigrations are the hand-written migrations from before the switch to Atlas.
// They are frozen: new schema changes are generated with Atlas into sql/ (see README).
func legacyMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		migrate20200000000000(),
		migrate20210709193000(),
		migrate20220126230000(),
		migrate20220713000000(),
		migrate20221119131300(),
		migrate20221210000000(),
		migrate20230825000000(),
		migrate20230904000000(),
		migrate20230530000000(),
		migrate20230904100000(),
		migrate20230826000000(),
		migrate20231003000000(),
		migrate20231015000000(),
		migrate20231023000000(),
		migrate20240101000000(),
		migrate20240102000000(),
		migrate20240103000000(),
		migrate20240207000000(),
		migrate20240311000000(),
		migrate20240312000000(),
		migrate20240316000000(),
		migrate20240317000000(),
		migrate20240318000000(),
		migrate20240319000000(),
		migrate20240320000000(),
		migrate20240327000000(),
		migrate20240402000000(),
		migrate20240405000000(),
		migrate20240510000000(),
		migrate20240511000000(),
		migrate20240512000000(),
		migrate20240706000000(),
		migrate20240824000000(),
		migrate20241023000000(),
		migrate20250412000000(),
	}
}

// postBaselineMigrations is frozen too; these run before the post-baseline SQL files.
// Keep their IDs for existing databases, but add all new migrations in sql/.
func postBaselineMigrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		migrate20261105000000(),
		migrate20261105010000(),
		migrate20261105020000(),
	}
}

// TableOptions are the options every table is created with.
// Without them, the `COLLATE` will not be set correctly.
const TableOptions = "ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci"

// Migrate brings the database to the latest schema by running the legacy migrations
// followed by the Atlas-generated SQL migrations. Applied migrations are tracked in the
// `migrations` table.
func Migrate(db *gorm.DB) error {
	sqlMigrations, err := sqlMigrations()
	if err != nil {
		return err
	}
	db = db.Set("gorm:table_options", TableOptions)
	log.Info("starting migration")
	start := time.Now()
	migrations := append(legacyMigrations(), postBaselineMigrations()...)
	migrations = append(migrations, sqlMigrations...)
	err = gormigrate.New(db, &gormigrate.Options{
		TableName:                 gormigrate.DefaultOptions.TableName,
		IDColumnName:              gormigrate.DefaultOptions.IDColumnName,
		IDColumnSize:              gormigrate.DefaultOptions.IDColumnSize,
		UseTransaction:            true,
		ValidateUnknownMigrations: true,
	}, migrations).Migrate()
	log.WithField("elapsed", time.Since(start)).Info("migration done")
	return err
}

// MigrateFromBaseline creates the latest schema on an empty database from the Atlas
// baseline snapshot plus the SQL migrations after it.
//
// It produces the same schema as Migrate, but skips the legacy migrations (and the seed
// data they insert) and does not record anything in the `migrations` table.
// It is much faster than Migrate and meant for tests.
func MigrateFromBaseline(db *gorm.DB) error {
	files, err := sqlMigrationFiles()
	if err != nil {
		return err
	}
	if err := files[0].apply(db); err != nil {
		return fmt.Errorf("applying %s: %w", files[0].name, err)
	}
	for _, migration := range postBaselineMigrations() {
		if err := migration.Migrate(db); err != nil {
			return fmt.Errorf("applying %s: %w", migration.ID, err)
		}
	}
	for _, f := range files[1:] {
		if err := f.apply(db); err != nil {
			return fmt.Errorf("applying %s: %w", f.name, err)
		}
	}
	return nil
}
