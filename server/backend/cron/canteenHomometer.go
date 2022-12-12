package cron

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AccessPoint struct {
	Target     string      `json:"target"`
	DataPoints [][]float32 `json:"datapoints"`
}

type CanteenApInformation struct {
	CanteenId string
	Target    string
	MaxCount  uint32
}

var (
	CANTEENS = []CanteenApInformation{
		{CanteenId: "mensa-arcisstr",
			Target:   "ap.ap*-?bn*.ssid.*",
			MaxCount: 450},
		{CanteenId: "mensa-garching",
			Target:   "ap.apa{01,02,03,04,05,06,07,08,09,10,11,12,13}*-?mg*.ssid.*",
			MaxCount: 1000},
		{CanteenId: "mensa-leopoldstr",
			Target:   "ap.ap*-?lm*.ssid.*",
			MaxCount: 500},
		{CanteenId: "mensa-lothstr",
			Target:   "ap.{apa06-0rh,apa05-0rh,apa03-1rh,apa02-1rh}.ssid.*",
			MaxCount: 110},
		{CanteenId: "mensa-martinsried",
			Target:   "ap.ap*-?ij*.ssid.*",
			MaxCount: 100},
		// Same as for the stucafe-parsing since I can not distinguish them
		{CanteenId: "mensa-pasing",
			Target:   "ap.apa{15,18}-?rl*.ssid.*",
			MaxCount: 50},
		{CanteenId: "mensa-weihenstephan",
			Target:   "ap.ap*-?qz*.ssid.*",
			MaxCount: 250},
		// No data found. 'Kantine' has only a few users connected and has a wrong pattern of connected users: http://wlan.lrz.de/apstat/filter/Unterbezirk/fs/
		{CanteenId: "stubistro-arcisstr",
			Target:   "ap.apa01-kfs.ssid.*",
			MaxCount: 25},
		// No data found
		{CanteenId: "stubistro-goethestr",
			Target:   "",
			MaxCount: 1000},
		// No data found. 'Mensaria' has only a few users connected: http://wlan.lrz.de/apstat/filter/Unterbezirk/if/
		{CanteenId: "stubistro-grosshadern",
			Target:   "",
			MaxCount: 1000},
		// No data found
		{CanteenId: "stubistro-rosenheim",
			Target:   "",
			MaxCount: 1000},
		// No data found
		{CanteenId: "stubistro-schellingstr",
			Target:   "",
			MaxCount: 1000},
		// No data found
		{CanteenId: "stubistro-adalbertstr",
			Target:   "",
			MaxCount: 1000},
		{CanteenId: "stubistro-martinsried",
			Target:   "ap.ap*-?iv*.ssid.*",
			MaxCount: 125},
		// No data found
		{CanteenId: "stucafe-akademie-weihenstephan",
			Target:   "",
			MaxCount: 1000},
		{CanteenId: "stucafe-boltzmannstr",
			Target:   "ap.apa{14,15,16}-?mg*.ssid.*",
			MaxCount: 1000},
		// No APs for the stucafe found: http://wlan.lrz.de/apstat/filter/Unterbezirk/ef/
		{CanteenId: "stucafe-connollystr",
			Target:   "",
			MaxCount: 1000},
		{CanteenId: "stucafe-garching",
			Target:   "ap.apa{14,15,16}*-?mg*.ssid.*",
			MaxCount: 1000},
		// No data found: http://wlan.lrz.de/apstat/filter/Unterbezirk/rf/
		{CanteenId: "stucafe-karlstr",
			Target:   "",
			MaxCount: 1000},
		// Same as for the mensa-parsing since I can not distinguish them
		{CanteenId: "stucafe-pasing",
			Target:   "ap.apa{15,18}-?rl*.ssid.*",
			MaxCount: 50},
		// No data found
		{CanteenId: "ipp-bistro",
			Target:   "",
			MaxCount: 1000},
		// No data available since the RBG does not provide stats
		{CanteenId: "fmi-bistro",
			Target:   "",
			MaxCount: 1000},
		// No data found
		{CanteenId: "mediziner-mensa",
			Target:   "",
			MaxCount: 1000},
		// Don't know which APs are for the canteen: http://wlan.lrz.de/apstat/filter/Unterbezirk/ua/
		{CanteenId: "mensa-straubing",
			Target:   "",
			MaxCount: 1000},
	}
)

/*
BASE_URL is the base URL for the required format.
Contains the 'XXX' placeholder that has to replaced with the Target property of the
CanteenApInformation when performing a request.
*/
const BASE_URL = "http://graphite-kom.srv.lrz.de/render/?from=-10min&target=XXX&format=json"

func (c *CronService) canteenHomometerCron() error {
	log.Info("Updating canteen homometer stats...")
	for _, canteen := range CANTEENS {
		if len(canteen.Target) <= 0 {
			log.Debug("Skipping canteen homometer stats for '", canteen.CanteenId, "', since there is no target.")
			continue
		}

		log.Debug("Updating canteen homometer stats for: ", canteen.CanteenId)
		aps := requestApData(&canteen)
		if len(aps) <= 0 {
			log.Debug("No canteen homometer data points found for: ", canteen.CanteenId)
			continue
		}

		count := sumApCounts(aps)
		updateDb(&canteen, count, c.db)
		log.Debug("Canteen homometer stats (", count, ") updated for: ", canteen.CanteenId)
	}
	log.Info("Canteen homometer stats updated.")
	return nil
}

func sumApCounts(aps []AccessPoint) uint32 {
	total := uint32(0)
	for _, ap := range aps {
		_, count, err := parseAp(&ap)
		if err != nil {
			log.WithError(err).Error("Canteen homometer getting the count failed for access point: ", ap.Target)
			continue
		}
		total += uint32(count)
	}
	return total
}

func updateDb(canteen *CanteenApInformation, count uint32, db *gorm.DB) error {
	entry := model.CanteenHomometer{
		CanteenId: canteen.CanteenId,
		Count:     count,
		MaxCount:  canteen.MaxCount,
		Percent:   (float32(count) / float32(canteen.MaxCount)) * 100,
		Timestamp: time.Now(),
	}

	res := db.Model(&model.CanteenHomometer{}).Where(model.CanteenHomometer{CanteenId: canteen.CanteenId}).Updates(&entry)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return db.Create(&entry).Error
	}
	return nil
}

func requestApData(canteen *CanteenApInformation) []AccessPoint {
	// Perform web request
	url := strings.Replace(BASE_URL, "XXX", canteen.Target, 1)
	resp, err := http.Get(url)
	if err != nil {
		log.WithError(err).Error("Canteen homometer web request failed for: ", canteen.CanteenId)
		return []AccessPoint{}
	}

	// Ensure we close the body once we leave this function
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	// Read the body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Canteen homometer reading body failed for: ", canteen.CanteenId)
		return []AccessPoint{}
	}

	// Parse as JSON
	aps := []AccessPoint{}
	err = json.Unmarshal(body, &aps)
	if err != nil {
		log.WithError(err).Error("Canteen homometer parsing output as JSON failed for: ", canteen.CanteenId)
		log.Trace("Body for '", canteen.CanteenId, "': ", body)
		return []AccessPoint{}
	}

	return aps
}

func parseAp(ap *AccessPoint) (string, uint32, error) {
	parts := strings.Split(ap.Target, ".")
	if len(parts) < 2 {
		return "", 0, errors.New("invalid access point name")
	}

	if len(ap.DataPoints) <= 0 {
		return ap.Target, 0, nil
	}

	// Check the last data point for the number of users
	count := uint32(0)
	lastTime := float32(0.0)
	for _, dp := range ap.DataPoints {
		if len(dp) == 2 && lastTime < dp[1] {
			lastTime = dp[1]
			count = uint32(dp[0])
		}
	}
	return ap.Target, count, nil
}
