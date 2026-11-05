// Command atlasloader prints the desired MySQL schema for Atlas (https://atlasgo.io).
//
// The desired schema consists of the GORM models plus legacy.sql, which declares the
// tables that exist in production but have no GORM model, and the foreign keys whose
// production names GORM cannot generate.
// It is used by the external_schema data source in server/atlas.hcl.
package main

import (
	_ "embed"
	"fmt"
	"os"
	"regexp"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/TUM-Dev/Campus-Backend/server/backend/migration"
)

//go:embed legacy.sql
var legacySchema string

// gormschema ignores gorm:table_options, so the options migration.Migrate sets are appended here.
var createTableEnd = regexp.MustCompile(`(?m)^(CREATE TABLE .*\));$`)

func main() {
	stmts, err := gormschema.New("mysql").Load(migration.Models()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	stmts = createTableEnd.ReplaceAllString(stmts, "$1 "+migration.TableOptions+";")
	fmt.Println(stmts)
	fmt.Println(legacySchema)
}
