package dbservice

import (
	"log"
)

func UpdateFuture(day int, hour int, avg int, apName string) {
	query := ""
	switch hour {
	case 0:
		query = "UPDATE future SET T0 = ? WHERE AP_NAME = ? AND Day = ?"
	case 1:
		query = "UPDATE future SET T1 = ? WHERE AP_NAME = ? AND Day = ?"
	case 2:
		query = "UPDATE future SET T2 = ? WHERE AP_NAME = ? AND Day = ?"
	case 3:
		query = "UPDATE future SET T3 = ? WHERE AP_NAME = ? AND Day = ?"
	case 4:
		query = "UPDATE future SET T4 = ? WHERE AP_NAME = ? AND Day = ?"
	case 5:
		query = "UPDATE future SET T5 = ? WHERE AP_NAME = ? AND Day = ?"
	case 6:
		query = "UPDATE future SET T6 = ? WHERE AP_NAME = ? AND Day = ?"
	case 7:
		query = "UPDATE future SET T7 = ? WHERE AP_NAME = ? AND Day = ?"
	case 8:
		query = "UPDATE future SET T8 = ? WHERE AP_NAME = ? AND Day = ?"
	case 9:
		query = "UPDATE future SET T9 = ? WHERE AP_NAME = ? AND Day = ?"
	case 10:
		query = "UPDATE future SET T10 = ? WHERE AP_NAME = ? AND Day = ?"
	case 11:
		query = "UPDATE future SET T11 = ? WHERE AP_NAME = ? AND Day = ?"
	case 12:
		query = "UPDATE future SET T12 = ? WHERE AP_NAME = ? AND Day = ?"
	case 13:
		query = "UPDATE future SET T13 = ? WHERE AP_NAME = ? AND Day = ?"
	case 14:
		query = "UPDATE future SET T14 = ? WHERE AP_NAME = ? AND Day = ?"
	case 15:
		query = "UPDATE future SET T15 = ? WHERE AP_NAME = ? AND Day = ?"
	case 16:
		query = "UPDATE future SET T16 = ? WHERE AP_NAME = ? AND Day = ?"
	case 17:
		query = "UPDATE future SET T17 = ? WHERE AP_NAME = ? AND Day = ?"
	case 18:
		query = "UPDATE future SET T18 = ? WHERE AP_NAME = ? AND Day = ?"
	case 19:
		query = "UPDATE future SET T19 = ? WHERE AP_NAME = ? AND Day = ?"
	case 20:
		query = "UPDATE future SET T20 = ? WHERE AP_NAME = ? AND Day = ?"
	case 21:
		query = "UPDATE future SET T21 = ? WHERE AP_NAME = ? AND Day = ?"
	case 22:
		query = "UPDATE future SET T22 = ? WHERE AP_NAME = ? AND Day = ?"
	case 23:
		query = "UPDATE future SET T23 = ? WHERE AP_NAME = ? AND Day = ?"
	default:
		log.Printf("Hour should be  >= 0 and < 24, but was: %d", hour)
		return
	}
	runQuery(query, avg, apName, day)
}

// Populates the future table with 15 name entries
// for each access point (one entry per day).
func PopulateFutureTable(accessPoints []AccessPoint) {
	db := InitDB(heatmapDB)
	for _, ap := range accessPoints {
		for day := 0; day < 15; day++ {
			for hour := 0; hour < 24; hour++ {
				InsertFuture(ap.Name, day, hour)
			}
		}
	}
	db.Close()
}

func InsertFuture(apName string, day int, hour int) {
	insertQuery := "INSERT INTO future(AP_Name, Day) VALUES (?, ?)"
	runQuery(insertQuery, apName, day)
}

func UpdateTomorrow() {
	deleteDayEq30 := "DELETE FROM future WHERE Day = 0"
	runQuery(deleteDayEq30)

	decreaseDays := "UPDATE future SET Day = Day - 1"
	runQuery(decreaseDays)
}

func GetPredictionForSingleAP(day int, hour int, name string) AccessPoint {
	db := InitDB(heatmapDB)
	query := ""

	switch hour {
	case 0:
		query = "SELECT T0 FROM future WHERE AP_Name = ? AND Day = ?"
	case 1:
		query = "SELECT T1 FROM future WHERE AP_Name = ? AND Day = ?"
	case 2:
		query = "SELECT T2 FROM future WHERE AP_Name = ? AND Day = ?"
	case 3:
		query = "SELECT T3 FROM future WHERE AP_Name = ? AND Day = ?"
	case 4:
		query = "SELECT T4 FROM future WHERE AP_Name = ? AND Day = ?"
	case 5:
		query = "SELECT T5 FROM future WHERE AP_Name = ? AND Day = ?"
	case 6:
		query = "SELECT T6 FROM future WHERE AP_Name = ? AND Day = ?"
	case 7:
		query = "SELECT T7 FROM future WHERE AP_Name = ? AND Day = ?"
	case 8:
		query = "SELECT T8 FROM future WHERE AP_Name = ? AND Day = ?"
	case 9:
		query = "SELECT T9 FROM future WHERE AP_Name = ? AND Day = ?"
	case 10:
		query = "SELECT T10 FROM future WHERE AP_Name = ? AND Day = ?"
	case 11:
		query = "SELECT T11 FROM future WHERE AP_Name = ? AND Day = ?"
	case 12:
		query = "SELECT T12 FROM future WHERE AP_Name = ? AND Day = ?"
	case 13:
		query = "SELECT T13 FROM future WHERE AP_Name = ? AND Day = ?"
	case 14:
		query = "SELECT T14 FROM future WHERE AP_Name = ? AND Day = ?"
	case 15:
		query = "SELECT T15 FROM future WHERE AP_Name = ? AND Day = ?"
	case 16:
		query = "SELECT T16 FROM future WHERE AP_Name = ? AND Day = ?"
	case 17:
		query = "SELECT T17 FROM future WHERE AP_Name = ? AND Day = ?"
	case 18:
		query = "SELECT T18 FROM future WHERE AP_Name = ? AND Day = ?"
	case 19:
		query = "SELECT T19 FROM future WHERE AP_Name = ? AND Day = ?"
	case 20:
		query = "SELECT T20 FROM future WHERE AP_Name = ? AND Day = ?"
	case 21:
		query = "SELECT T21 FROM future WHERE AP_Name = ? AND Day = ?"
	case 22:
		query = "SELECT T22 FROM future WHERE AP_Name = ? AND Day = ?"
	case 23:
		query = "SELECT T23 FROM future WHERE AP_Name = ? AND Day = ?"
	default:
		log.Printf("Hour should be  >= 0 and < 24, but was: %d", hour)
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

func GetFutureForAllAPs(day int, hour int) []AccessPoint {
	db := InitDB(heatmapDB)
	query := "SELECT DISTINCT AP_Name, Load, Max, Min FROM future WHERE Day = ? AND Hour + ?"

	rows, err := db.Query(query, day, hour)
	if err != nil {
		log.Println(err)
		return []AccessPoint{}
	}
	defer rows.Close()

	var apList []AccessPoint
	apList = make([]AccessPoint, 0)

	for rows.Next() {
		var ap AccessPoint
		err = rows.Scan(&ap.Name, &ap.Load, &ap.Max, &ap.Min)
		if err == nil {
			apList = append(apList, ap)
		}
	}

	db.Close()
	return apList
}
