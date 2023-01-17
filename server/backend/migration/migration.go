// Package migration contains functions related to database changes and executes them
package migration

import (
	"errors"
	"fmt"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
)

// TumDBMigrator contains a reference to our database
type TumDBMigrator struct {
	database          *gorm.DB
	shouldAutoMigrate bool
}

// New creates a new TumDBMigrator with a database
func New(db *gorm.DB, shouldAutoMigrate bool) TumDBMigrator {
	return TumDBMigrator{database: db, shouldAutoMigrate: shouldAutoMigrate}
}

// Migrate starts the migration either by using AutoMigrate in development environments or manually in prod
func (m TumDBMigrator) Migrate() error {
	if m.shouldAutoMigrate {
		log.Info("Using automigration")
		err := m.database.AutoMigrate(
			&model.TopNews{},
			&model.Crontab{},
			&model.Files{},
			&model.NewsSource{},
			&model.NewsAlert{},
			&model.News{},
			&model.CanteenHeadCount{},
		)
		return err
	}
	log.Info("Using manual migration")
	mig := gormigrate.New(m.database, gormigrate.DefaultOptions, []*gormigrate.Migration{
		m.migrate20210709193000(),
		m.migrate20220126230000(),
		m.migrate20220713000000(),
		m.migrate20221119131300(),
		m.migrate20221210000000(),
	})
	err := mig.Migrate()
	return err

}

func SafeEnumMigrate(tx *gorm.DB, table interface{}, column string, additionalTypes ...string) error {
	enumTypes, err := getEnumTypesFromDB(tx, table, column)

	if err != nil {
		return err
	}

	enumTypes = append(enumTypes, additionalTypes...)
	enumTypes = ensureUnique(enumTypes)

	return alterEnumColumn(tx, table, column, enumTypes)
}

func SafeEnumRollback(tx *gorm.DB, table interface{}, column string, rollbackTypes ...string) error {
	enumTypes, err := getEnumTypesFromDB(tx, table, column)

	if err != nil {
		return err
	}

	enumTypes = removeTypes(enumTypes, rollbackTypes...)

	return alterEnumColumn(tx, table, column, enumTypes)
}

func ensureUnique(types []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range types {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func alterEnumColumn(tx *gorm.DB, table interface{}, column string, types []string) error {
	enum := buildEnum(types)

	stmt := &gorm.Statement{DB: tx}
	stmt.Parse(&table)
	tableName := stmt.Schema.Table

	rawQuery := fmt.Sprintf(
		"ALTER TABLE %s MODIFY %s %s;",
		tableName,
		column,
		enum,
	)

	tx = tx.Exec(rawQuery)

	if tx.Error != nil {
		return errors.New("could not alter table")
	}

	return nil
}

func getEnumTypesFromDB(tx *gorm.DB, table interface{}, column string) ([]string, error) {
	columnType, err := tx.Migrator().ColumnTypes(&table)

	if err != nil {
		return nil, errors.New("could not get column types")
	}

	enumTypes, err := getEnumTypes(columnType, column)

	if err != nil {
		return nil, err
	}

	return enumTypes, nil
}

func removeTypes(types []string, rollbackTypes ...string) []string {
	for _, t := range rollbackTypes {
		for i, tt := range types {
			if tt == t {
				types = append(types[:i], types[i+1:]...)
			}
		}
	}

	return types
}

func getEnumTypes(columTypes []gorm.ColumnType, column string) ([]string, error) {
	var cType string

	for _, t := range columTypes {
		if t.Name() == column {
			if t, ok := t.ColumnType(); ok {
				cType = t
				break
			} else {
				return nil, errors.New("could not get column type")
			}
		}
	}

	if !strings.Contains(strings.ToLower(cType), "enum") {
		return nil, errors.New("column is not an enum")
	}

	leftTrimmed := strings.TrimLeft(cType, "enum")
	trimmed := strings.Trim(leftTrimmed, "()")

	splitted := strings.Split(trimmed, ",")

	types := make([]string, len(splitted))

	for i, s := range splitted {
		types[i] = strings.TrimSpace(s)
		types[i] = strings.Trim(types[i], "''")
	}

	return types, nil
}

func buildEnum(types []string) string {
	str := "enum("

	for _, t := range types {
		str += fmt.Sprintf("'%s',", t)
	}

	str = strings.TrimRight(str, ",")

	str += ")"

	return str
}
