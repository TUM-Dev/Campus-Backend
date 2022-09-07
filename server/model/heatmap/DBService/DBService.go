package DBService

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	//path to the database
	heatmapDB = "./data/sqlite/heatmap.db"
)

var DB = InitDB(heatmapDB)

type AccessPoint struct {
	ID      string // primary key
	Address string
	Room    string
	Name    string
	Floor   string
	Status  string
	Type    string
	Load    string
	Lat     string
	Long    string
	Max     int
	Min     int
}

// Opens the database and returns a non-nil pointer if successfull
func InitDB(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Panicf("Could not open the database! %v", err)
	}

	if db == nil {
		panic("DB pointer is nil!")
	}
	return db
}

func RetrieveAPsOfTUM(withCoordinate bool) []AccessPoint {
	var query string

	if withCoordinate {
		query = `
			SELECT ID, Address, Room, Name, Floor, Load, Lat, Long
			FROM apstat
			WHERE Address LIKE '%TUM%'
			AND Lat!='lat'
			AND Long!='long'
		`
	} else {
		query = `
			SELECT ID, Address, Room, Name, Floor, Load, Lat, Long
			FROM apstat
			WHERE Address LIKE '%TUM%'
			AND Lat='lat'
			AND Long='long'
		`
	}

	return RetrieveAPsFromTUM(query)
}

// Queries the 'apstat' table and
// returns all rows where 'address' contains "TUM" and
// the columns 'Lat', 'Long' are yet unassigned
func RetrieveAPsFromTUM(query string) []AccessPoint {
	db := InitDB(heatmapDB)

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []AccessPoint
	for rows.Next() {
		item := AccessPoint{}
		err := rows.Scan(&item.ID, &item.Address, &item.Room, &item.Name, &item.Floor, &item.Load, &item.Lat, &item.Long)

		if err != nil {
			panic(err)
		}
		result = append(result, item)
	}

	db.Close()
	return result
}

func UpdateLatLong(column, newValue, id string) {
	query := ""
	if column == "Lat" {
		query = "UPDATE apstat SET Lat = ? WHERE ID = ?"	
	} else {
		query = "UPDATE apstat SET Long = ? WHERE ID = ?"
	}
	runQuery(query, newValue, id)
}

func UpdateMinMax(column, apName string, newValue int) {
	query := ""
	if column == "Max" {
		query = "UPDATE apstat SET Max = ? WHERE Name = ?"	
	} else {
		query = "UPDATE apstat SET Min = ? WHERE Name = ?"
	}
	runQuery(query, newValue, apName)
}

// Updates name of 'currName' column to 'newName'.
func UpdateColumnName(tableName, currName, newName string) {
	query := fmt.Sprintf(`ALTER TABLE %s RENAME COLUMN %s TO %s;`, tableName, currName, newName)
	runQuery(query)
}

// Prepares sqlite 'query' and executes it with optional 'params'.
func runQuery(query string, params ...interface{}) {
	stmt, err := DB.Prepare(query)

	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		panic(err)
	}
}
