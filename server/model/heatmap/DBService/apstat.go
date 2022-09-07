package DBService

import (
	"log"
)

func GetAllNames() []string {
	query := `
	SELECT Name
	FROM apstat
	WHERE Address LIKE '%TUM%'
		AND Lat!='lat'
		AND Long!='long'
	`
	rows, err := DB.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)

		if err != nil {
			panic(err)
		}
		names = append(names, name)
	}
	return names
}

func GetAccessPointByName(name string) *AccessPoint {
	stmt :=`
	SELECT Name, Lat, Long, Load
	FROM apstat
	WHERE Name=?
	`

	row := DB.QueryRow(stmt, name)
	result := AccessPoint{}
	switch err := row.Scan(&result.Name, &result.Lat, &result.Long, &result.Load); err {
	case nil:
		log.Printf("Returning result for access point with name %s!", name)
	default:
		log.Printf("No access point found with name %s in the DB!", name)
		log.Println("Returning empty result.")
	}
	return &result
}
