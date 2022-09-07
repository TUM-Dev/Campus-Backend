package DBService

import (
	"database/sql"
	"log"
)

func UpdateHistory(day int, hour int, avg int, apName string) {
	query := ""
	switch hour {
	case 0:
		query = "UPDATE history SET T0 = ? WHERE Day = ?"
	case 1:
		query = "UPDATE history SET T1 = ? WHERE Day = ?"
	case 2:
		query = "UPDATE history SET T2 = ? WHERE Day = ?"
	case 3:
		query = "UPDATE history SET T3 = ? WHERE Day = ?"
	case 4:
		query = "UPDATE history SET T4 = ? WHERE Day = ?"
	case 5:
		query = "UPDATE history SET T5 = ? WHERE Day = ?"
	case 6:
		query = "UPDATE history SET T6 = ? WHERE Day = ?"
	case 7:
		query = "UPDATE history SET T7 = ? WHERE Day = ?"
	case 8:
		query = "UPDATE history SET T8 = ? WHERE Day = ?"
	case 9:
		query = "UPDATE history SET T9 = ? WHERE Day = ?"
	case 10:
		query = "UPDATE history SET T10 = ? WHERE Day = ?"
	case 11:
		query = "UPDATE history SET T11 = ? WHERE Day = ?"
	case 12:
		query = "UPDATE history SET T12 = ? WHERE Day = ?"
	case 13:
		query = "UPDATE history SET T13 = ? WHERE Day = ?"
	case 14:
		query = "UPDATE history SET T14 = ? WHERE Day = ?"
	case 15:
		query = "UPDATE history SET T15 = ? WHERE Day = ?"
	case 16:
		query = "UPDATE history SET T16 = ? WHERE Day = ?"
	case 17:
		query = "UPDATE history SET T17 = ? WHERE Day = ?"
	case 18:
		query = "UPDATE history SET T18 = ? WHERE Day = ?"
	case 19:
		query = "UPDATE history SET T19 = ? WHERE Day = ?"
	case 20:
		query = "UPDATE history SET T20 = ? WHERE Day = ?"
	case 21:
		query = "UPDATE history SET T21 = ? WHERE Day = ?"
	case 22:
		query = "UPDATE history SET T22 = ? WHERE Day = ?"
	case 23:
		query = "UPDATE history SET T23 = ? WHERE Day = ?"
	default:
		log.Panicf("Hour should be  >= 0 and < 24, but was: %d", hour)
	}
	runQuery(query, avg, apName, day)
}

// Populates the last31days table with 31 name entries
// for each access point (one entry per day).
func PopulateHistoryTable(accessPoints []AccessPoint) {
	db := InitDB(heatmapDB)
	for _, ap := range accessPoints {
		for day := 0; day < 31; day++ {
			for hour := 0; hour < 24; hour++ {
				InsertHistory(ap.Name, day, hour, db)
			}
		}
	}
	db.Close()
}

func InsertHistory(apName string, day int, hour int, db *sql.DB) {
	insertQuery := "INSERT INTO history(AP_Name, Day, Hour) VALUES (?, ?, ?)"
	runQuery(insertQuery, apName, day, hour)
}

// returns a map of AP names, whose last 30 day data
// hasn't been stored in the database yet
func GetUnprocessedAPs() map[string]bool {
	db := InitDB(heatmapDB)
	query := `
		SELECT DISTINCT AP_Name
		FROM history
		WHERE T0 IS NULL
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return make(map[string]bool)
	}
	defer rows.Close()

	var names = make(map[string]bool)
	for rows.Next() {
		var ap string
		rows.Scan(&ap)
		names[ap] = true
	}

	db.Close()
	return names
}

// queries the database and returns the intensity of
// the access point, based on the selected day and hour
func GetHistoryForSingleAP(name string, day int, hour int) AccessPoint {
	db := InitDB(heatmapDB)
	query := ""

	switch hour {
	case 0:
		query = "SELECT T0 FROM history WHERE Day = ?"
	case 1:
		query = "SELECT T1 FROM history WHERE Day = ?"
	case 2:
		query = "SELECT T2 FROM history WHERE Day = ?"
	case 3:
		query = "SELECT T3 FROM history WHERE Day = ?"
	case 4:
		query = "SELECT T4 FROM history WHERE Day = ?"
	case 5:
		query = "SELECT T5 FROM history WHERE Day = ?"
	case 6:
		query = "SELECT T6 FROM history WHERE Day = ?"
	case 7:
		query = "SELECT T7 FROM history WHERE Day = ?"
	case 8:
		query = "SELECT T8 FROM history WHERE Day = ?"
	case 9:
		query = "SELECT T9 FROM history WHERE Day = ?"
	case 10:
		query = "SELECT T10 FROM history WHERE Day = ?"
	case 11:
		query = "SELECT T11 FROM history WHERE Day = ?"
	case 12:
		query = "SELECT T12 FROM history WHERE Day = ?"
	case 13:
		query = "SELECT T13 FROM history WHERE Day = ?"
	case 14:
		query = "SELECT T14 FROM history WHERE Day = ?"
	case 15:
		query = "SELECT T15 FROM history WHERE Day = ?"
	case 16:
		query = "SELECT T16 FROM history WHERE Day = ?"
	case 17:
		query = "SELECT T17 FROM history WHERE Day = ?"
	case 18:
		query = "SELECT T18 FROM history WHERE Day = ?"
	case 19:
		query = "SELECT T19 FROM history WHERE Day = ?"
	case 20:
		query = "SELECT T20 FROM history WHERE Day = ?"
	case 21:
		query = "SELECT T21 FROM history WHERE Day = ?"
	case 22:
		query = "SELECT T22 FROM history WHERE Day = ?"
	case 23:
		query = "SELECT T23 FROM history WHERE Day = ?"
	default:
		log.Panicf("Hour should be  >= 0 and < 24, but was: %d", hour)
	}

	row := db.QueryRow(query, name, day)
	result := AccessPoint{}
	switch err := row.Scan(&result.Load); err {
	case nil:
		log.Println("Returning result from history!")
	default:
		log.Println("No data found in history! Returning empty result.")
	}
	
	db.Close()
	return result
}

// queries the database and returns a list of the names and intensities (network load)
// of all access points, based on the selected day and hour.
func GetHistoryForAllAPs(day int, hour int) []AccessPoint {
	db := InitDB(heatmapDB)
	query := ""

	switch hour {
	case 0:
		query = "SELECT DISTINCT AP_Name, T0, Max, Min FROM history WHERE Day = ?"
	case 1:
		query = "SELECT DISTINCT AP_Name, T1, Max, Min FROM history WHERE Day = ?"
	case 2:
		query = "SELECT DISTINCT AP_Name, T2, Max, Min FROM history WHERE Day = ?"
	case 3:
		query = "SELECT DISTINCT AP_Name, T3, Max, Min FROM history WHERE Day = ?"
	case 4:
		query = "SELECT DISTINCT AP_Name, T4, Max, Min FROM history WHERE Day = ?"
	case 5:
		query = "SELECT DISTINCT AP_Name, T5, Max, Min FROM history WHERE Day = ?"
	case 6:
		query = "SELECT DISTINCT AP_Name, T6, Max, Min FROM history WHERE Day = ?"
	case 7:
		query = "SELECT DISTINCT AP_Name, T7, Max, Min FROM history WHERE Day = ?"
	case 8:
		query = "SELECT DISTINCT AP_Name, T8, Max, Min FROM history WHERE Day = ?"
	case 9:
		query = "SELECT DISTINCT AP_Name, T9, Max, Min FROM history WHERE Day = ?"
	case 10:
		query = "SELECT DISTINCT AP_Name, T10, Max, Min FROM history WHERE Day = ?"
	case 11:
		query = "SELECT DISTINCT AP_Name, T11, Max, Min FROM history WHERE Day = ?"
	case 12:
		query = "SELECT DISTINCT AP_Name, T12, Max, Min FROM history WHERE Day = ?"
	case 13:
		query = "SELECT DISTINCT AP_Name, T13, Max, Min FROM history WHERE Day = ?"
	case 14:
		query = "SELECT DISTINCT AP_Name, T14, Max, Min FROM history WHERE Day = ?"
	case 15:
		query = "SELECT DISTINCT AP_Name, T15, Max, Min FROM history WHERE Day = ?"
	case 16:
		query = "SELECT DISTINCT AP_Name, T16, Max, Min FROM history WHERE Day = ?"
	case 17:
		query = "SELECT DISTINCT AP_Name, T17, Max, Min FROM history WHERE Day = ?"
	case 18:
		query = "SELECT DISTINCT AP_Name, T18, Max, Min FROM history WHERE Day = ?"
	case 19:
		query = "SELECT DISTINCT AP_Name, T19, Max, Min FROM history WHERE Day = ?"
	case 20:
		query = "SELECT DISTINCT AP_Name, T20, Max, Min FROM history WHERE Day = ?"
	case 21:
		query = "SELECT DISTINCT AP_Name, T21, Max, Min FROM history WHERE Day = ?"
	case 22:
		query = "SELECT DISTINCT AP_Name, T22, Max, Min FROM history WHERE Day = ?"
	case 23:
		query = "SELECT DISTINCT AP_Name, T23, Max, Min FROM history WHERE Day = ?"
	default:
		log.Panicf("Hour should be  >= 0 and < 24, but was: %d", hour)
	}
	
	rows, err := db.Query(query, day)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var apList []AccessPoint
	apList = make([]AccessPoint, 0)
	
	for rows.Next() {
		var ap AccessPoint
		err = rows.Scan(&ap.Name, &ap.Load, &ap.Max, &ap.Min);
		if err == nil {
			apList = append(apList, ap)
		}
	}
	
	db.Close()
	return apList
}

func UpdateToday(averages map[string][24]int) {
	// remove day 30 + add new day 0 = replace 30 with day -1 AND
	// shift days 0-29 by adding 1 day AND
	// update day 0 with new data
	
	deleteDayEq30 := "UPDATE history SET Day = -1 WHERE Day = 30"
	runQuery(deleteDayEq30)
	
	increaseDays := "UPDATE history SET Day = Day + 1"
	runQuery(increaseDays)

	for apName, dailyAvgs := range averages {
		for hour, dailyAvg := range dailyAvgs {
			// query := "UPDATE history SET Load = ? WHERE Day = 0"
			// runQuery(query, dailyAvg, apName)
			UpdateHistory(0, hour, dailyAvg, apName)
		}
	}
}