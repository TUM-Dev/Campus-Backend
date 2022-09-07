package DBService

import (
	"log"
)


// JOIN
func JoinMaxMin() {
	query := `
		UPDATE history
		SET Min = (SELECT Min
			FROM apstat WHERE history.AP_Name = Name),
		SET Max = (SELECT Max 
			FROM apstat WHERE history.AP_Name = Name)
	`
	runQuery(query)
}

func DeleteOldTables() {
	query := "DROP TABLE history"
	runQuery(query)
	query2 := "DROP TABLE future"
	runQuery(query2)
}

func TestExample() {
	for ap := 0; ap < 5; ap++ {
		for i := 0; i < 31; i++ {
			for j := 0; j < 24; j++ {
				UpdateHistory(i,j,100,"apa08-1w4")
				// log.Println("single update")
			}
		}
		log.Println("updated")
	}
}

func TestQuestionMark() {
	UpdateMinMax("T0","apa08-1w4", 99)
}