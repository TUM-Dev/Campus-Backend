package migration

import (
	"fmt"

	"github.com/go-gormigrate/gormigrate/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type tablesLiningToFiles struct {
	table          string
	field          string
	constraintName string
	typeDefinition string
}

// migrate20240318000000
// Enforced stricter constraints for files
// - introduced a unique index
// - made sure that all fields which should not be null are not null
// - Changed the primary key of files to be bigint
// - added previously missed foreign keys
func migrate20240318000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20240318000000",
		Migrate: func(tx *gorm.DB) error {
			// uniqueness
			if err := tx.Exec("alter table files add constraint url_unique unique (url)").Error; err != nil {
				return err
			}
			// nullability
			if err := tx.Exec("alter table files modify name text not null").Error; err != nil {
				return err
			}
			if err := tx.Exec("alter table files modify path text not null").Error; err != nil {
				return err
			}
			// remove FKs to files
			tablesWithFKsToFiles := []tablesLiningToFiles{
				{"news", "file", "news_ibfk_2", "bigint not null auto_increment"},
				{"newsSource", "icon", "newsSource_ibfk_1", "bigint not null auto_increment"},
				{"kino", "cover", "kino_ibfk_1", "bigint null auto_increment"},
				{"event", "file", "fkEventFile", "bigint null auto_increment"},
			}

			for _, t := range tablesWithFKsToFiles {
				if err := tx.Exec(fmt.Sprintf("alter table `%s` drop foreign key `%s`", t.table, t.constraintName)).Error; err != nil {
					return err
				}
			}
			// news_alert has an index instead of a FK
			tablesWithFKsToFiles = append(tablesWithFKsToFiles, tablesLiningToFiles{"news_alert", "news_alert", "news_alert", "bigint not null auto_increment"})
			if err := tx.Exec("alter table news_alert drop key FK_File").Error; err != nil {
				return err
			}
			if err := migrateField(tx, "news_alert", "file", "bigint not null"); err != nil {
				return err
			}
			// change PK
			if err := migrateField(tx, "files", "file", "bigint not null auto_increment"); err != nil {
				return err
			}
			// re-add FKs to files
			for _, t := range tablesWithFKsToFiles {
				if err := tx.Exec(fmt.Sprintf("alter table `%s` add constraint `%s` foreign key (`%s`) references files (file) on update cascade on delete cascade", t.table, t.constraintName, t.field)).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			log.Fatal("intentionally no rollback function as this would be lossy!")
			return nil
		},
	}
}
