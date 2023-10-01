package migration

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

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

	enumTypes = RemoveTypes(enumTypes, rollbackTypes...)

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
	enum := BuildEnum(types)

	stmt := &gorm.Statement{DB: tx}
	err := stmt.Parse(&table)
	if err != nil {
		return errors.New("could not parse enum table")
	}
	tableName := stmt.Schema.Table

	rawQuery := fmt.Sprintf(
		"ALTER TABLE %s MODIFY %s %s;",
		tableName,
		column,
		enum,
	)

	tx = tx.Exec(rawQuery)

	if tx.Error != nil {
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

	return EnumTypesFromString(cType)
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
	str := "enum("

	for _, t := range types {
		str += fmt.Sprintf("'%s',", t)
	}

	str = strings.TrimRight(str, ",")

	str += ")"

	return str
}
