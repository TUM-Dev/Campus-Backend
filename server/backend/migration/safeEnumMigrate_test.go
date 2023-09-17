package migration_test

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"testing"

	"github.com/TUM-Dev/Campus-Backend/server/backend/migration"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	TestEnumTypes = []string{
		"news", "mensa", "canteen", "iosNotifications",
		"iosActivityReset", "canteenHeadCount", "testType",
	}
)

func generateTestEnums() ([]string, [][]string) {
	var enums []string
	var enumTypes [][]string

	for i := 0; i < 10; i++ {
		enum, values := buildEnum()
		enums = append(enums, enum)
		enumTypes = append(enumTypes, values)
	}

	return enums, enumTypes
}

func buildEnum() (string, []string) {
	var enum []string
	for _, t := range TestEnumTypes {
		if rand.Intn(2) == 1 {
			enum = append(enum, t)
		}
	}

	if len(enum) == 0 {
		enum = append(enum, TestEnumTypes[rand.Intn(len(TestEnumTypes))])
	}

	return fmt.Sprintf("enum('%s')", strings.Join(enum, "', '")), enum
}

func isValidEnum(enum string) bool {
	match, err := regexp.MatchString("enum\\s*\\('\\w+'(,\\s*'\\w+')*\\);?", enum)

	if err != nil {
		log.WithError(err).Error("error matching regex")
		return false
	}

	return match
}

func TestEnumTypesFromString(t *testing.T) {
	enums, enumTypes := generateTestEnums()

	for i, enum := range enums {
		types, err := migration.EnumTypesFromString(enum)

		assert.Nil(t, err, "error should be nil")

		assert.Equalf(t, len(types), len(enumTypes[i]), "length of enum types does not match")

		for j, jType := range types {
			assert.Equalf(t, jType, enumTypes[i][j], "enum types do not match")
		}
	}
}

func TestBuildEnum(t *testing.T) {
	_, values := generateTestEnums()

	for _, enumTypes := range values {
		enum := migration.BuildEnum(enumTypes)

		assert.True(t, isValidEnum(enum))
	}
}

func TestRemoveTypes(t *testing.T) {
	types := TestEnumTypes

	removeTypes := []string{
		"news", "mensa", "canteen", "iosNotifications",
	}

	expectedTypes := []string{
		"iosActivityReset", "canteenHeadCount", "testType",
	}

	newArr := migration.RemoveTypes(types, removeTypes...)

	assert.Equal(t, len(newArr), len(expectedTypes), "length of new array does not match")

	assert.Equalf(t, newArr, expectedTypes, "new array does not match")
}
