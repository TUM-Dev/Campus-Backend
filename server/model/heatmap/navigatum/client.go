package navigatum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/TUM-Dev/Campus-Backend/model/heatmap/dbservice"
	"github.com/TUM-Dev/Campus-Backend/model/heatmap/RoomFinder"
)

const (
	//path to the database
	heatmapDB = "./data/sqlite/heatmap.db"
)

type Coords struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"lon"`
	Src  string  `json:"source"`
}

// JSON response from NavigaTUM
type NaviResponse struct {
	Type  string `json:"type"`
	Coord Coords `json:"coords"`
}

func ScrapeNavigaTUM(res RoomFinder.Result) (count int) {
	count = 0 // number of found coordinates

	for _, res := range res.Failures {
		var roomID string
		if strings.Contains(res.RoomNr, "OG") || res.RoomNr == "" || strings.Contains(res.RoomNr, "..") {
			roomID = res.BuildingNr
		} else {
			roomID = fmt.Sprintf("%s.%s", res.BuildingNr, res.RoomNr)
		}

		lat, long, found := getRoomCoordinates(roomID)

		if found {
			storeCoordinateOfAP(res, lat, long)
			count++
		} else { // if no exact coord was found, store that of the building
			lat, long, found = getRoomCoordinates(res.BuildingNr)
			if found {
				storeCoordinateOfAP(res, lat, long)
				count++
			}
		}
	}

	return
}

// stores the coordinates of the access point in the database
func storeCoordinateOfAP(result RoomFinder.Failure, lat, long string) {
	db := dbservice.InitDB(heatmapDB)
	dbservice.UpdateLatLong("Lat", lat, result.ID)
	dbservice.UpdateLatLong("Long", long, result.ID)
	db.Close()
}

// makes an HTTP GET request to nav.tum.sexy/api/get/{roomID}
// e.g. roomID := 5602.EG.001
func getRoomCoordinates(roomID string) (lat, long string, found bool) {
	lat, long, found = "", "", false

	url := fmt.Sprintf("https://nav.tum.sexy/api/get/%s", roomID)
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("GET request to %s failed!", url)
		return
	}

	if resp.StatusCode != 200 {
		log.Printf("%v for url: %s", resp.Status, url)
		return
	}
	defer resp.Body.Close()

	var respNavi NaviResponse
	err = json.NewDecoder(resp.Body).Decode(&respNavi)

	if err != nil {
		log.Printf("JSON decoding failed! %q", err)
		log.Printf("%v", resp.Body)
		return
	}

	lat = fmt.Sprintf("%f", respNavi.Coord.Lat)
	long = fmt.Sprintf("%f", respNavi.Coord.Long)
	found = true

	return
}
