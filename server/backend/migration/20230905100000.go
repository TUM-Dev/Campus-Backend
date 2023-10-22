package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Movie stores all movies
type Movie struct {
	Id    int32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;not null;"`
	Cover Files `gorm:"column:cover;type:int;not null"`
}

// Kino stores all movies
type Kino struct {
	Cover   Files  `gorm:"column:cover;type:int;not null"`
	Trailer string `gorm:"column:trailer;type:text;"`
}

// TableName sets the insert table name for this struct type
func (n *Kino) TableName() string {
	return "kino"
}

type Files struct{}

// TableName sets the insert table name for this struct type
func (f *Files) TableName() string {
	return "files"
}

// migrate20230905100000
// removes the unused trailer column
// makes the Cover FK into a not null field
// renames kino -> movie
// fixes the id being named kino
func (m TumDBMigrator) migrate20230905100000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "migrate20230905100000",
		Migrate: func(tx *gorm.DB) error {
			// fix the movie table
			if err := tx.Migrator().RenameTable(&Kino{}, &Movie{}); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&Movie{}, "trailer"); err != nil {
				return err
			}
			if err := tx.Migrator().RenameColumn(&Movie{}, "kino", "id"); err != nil {
				return err
			}
			if err := tx.Migrator().AlterColumn(&Movie{}, "cover"); err != nil {
				return err
			}
			return nil
		},

		Rollback: func(tx *gorm.DB) error {
			// rollback the kino table
			if err := tx.Migrator().RenameTable(&Kino{}, &Movie{}); err != nil {
				return err
			}
			if err := tx.Migrator().RenameColumn(&Movie{}, "id", "kino"); err != nil {
				return err
			}
			if err := tx.AutoMigrate(Kino{}); err != nil {
				return err
			}
			return nil
		},
	}
}
