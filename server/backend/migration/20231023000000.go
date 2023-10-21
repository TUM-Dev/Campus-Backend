package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

type KinoWithNullableFields struct {
	Year       null.String `gorm:"column:year;type:varchar(4)"`
	Runtime    null.String `gorm:"column:runtime;type:varchar(40)"`
	Genre      null.String `gorm:"column:genre;type:varchar(100)"`
	Director   null.String `gorm:"column:director;type:text"`
	Actors     null.String `gorm:"column:actors;type:text"`
	ImdbRating null.String `gorm:"column:rating;type:varchar(4)"`
	Location   null.String `gorm:"column:location;default:null"`
}

type KinoWithoutNullableFields struct {
	Year       string `gorm:"column:year;type:varchar(4);not null;"`
	Runtime    string `gorm:"column:runtime;type:varchar(40);not null;"`
	Genre      string `gorm:"column:genre;type:varchar(100);not null;"`
	Director   string `gorm:"column:director;type:text;not null;"`
	Actors     string `gorm:"column:actors;type:text;not null;"`
	ImdbRating string `gorm:"column:rating;type:varchar(4);not null;"`
}

// TableName sets the insert table name for this struct type
func (n *KinoWithNullableFields) TableName() string {
	return "kino"
}

// migrate20231023000000
// migrates the static data for the canteen rating system and adds the necessary cronjob entries
func (m TumDBMigrator) migrate20231023000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20231023000000",
		Migrate: func(tx *gorm.DB) error {
			return tx.Migrator().AutoMigrate(&KinoWithNullableFields{})
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec("DROP FROM kino WHERE year IS NULL OR runtime IS NULL OR genre IS NULL OR director IS NULL OR actors IS NULL OR rating IS NULL").Error; err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&KinoWithNullableFields{}, "location"); err != nil {
				return err
			}
			return tx.Migrator().AutoMigrate(&KinoWithoutNullableFields{})
		},
	}
}
