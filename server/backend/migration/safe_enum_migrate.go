package migration

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

func SafeEnumAdd(tx *gorm.DB, table interface{}, column string, additionalTypes ...string) error {
	enumTypes, err := getEnumTypesFromDB(tx, table, column)
	if err != nil {
		return err
	}

	return alterEnumColumn(tx, table, column, ensureUnique(append(enumTypes, additionalTypes...)))
}

func SafeEnumRemove(tx *gorm.DB, table interface{}, column string, rollbackTypes ...string) error {
	enumTypes, err := getEnumTypesFromDB(tx, table, column)
	if err != nil {
		return err
	}

	return alterEnumColumn(tx, table, column, RemoveTypes(enumTypes, rollbackTypes...))
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
	stmt := &gorm.Statement{DB: tx}
	err := stmt.Parse(&table)
	if err != nil {
		return errors.New("could not parse enum table")
	}

	if err := tx.Exec(fmt.Sprintf(
		"ALTER TABLE %s MODIFY %s %s;",
		stmt.Schema.Table,
		column,
		BuildEnum(types),
	)).Error; err != nil {
		log.WithError(err).Error("Error altering enum table")
		return errors.New("could not alter enum table")
	}
	return nil
}

func getEnumTypesFromDB(tx *gorm.DB, table interface{}, column string) ([]string, error) {
	columnType, err := tx.Migrator().ColumnTypes(&table)
	if err != nil {
		return nil, errors.New("could not get enum column types")
	}

	enumTypes, err := getEnumTypes(columnType, column)
	if err != nil {
		return nil, err
	}
	return enumTypes, nil
}

func RemoveTypes(types []string, rollbackTypes ...string) []string {
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
	for _, t := range columTypes {
		if t.Name() == column {
			if t, ok := t.ColumnType(); ok {
				return EnumTypesFromString(t)
			} else {
				return nil, errors.New("could not get column type")
			}
		}
	}
	return nil, errors.New("column does not exist")
}

func EnumTypesFromString(enum string) ([]string, error) {
	if !strings.Contains(strings.ToLower(enum), "enum") {
		return nil, errors.New("column is not an enum")
	}

	leftTrimmed := strings.TrimLeft(enum, "enum")
	rightTrimmed := strings.TrimRight(leftTrimmed, ";")
	trimmed := strings.Trim(rightTrimmed, "()")

	splitted := strings.Split(trimmed, ",")

	types := make([]string, len(splitted))

	for i, s := range splitted {
		types[i] = strings.TrimSpace(s)
		types[i] = strings.Trim(types[i], "'")
	}

	return types, nil
}

func BuildEnum(types []string) string {
	enums := strings.Join(types, "','")
	return fmt.Sprintf("enum('%s')", enums)
}
