package cron

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AccessPoint struct {
	Target     string      `json:"target"`
	DataPoints [][]float32 `json:"datapoints"`
}

type CanteenApInformation struct {
	CanteenId string
	Url       string
	MaxCount  uint32
}

var (
	CANTEENS = []CanteenApInformation{{CanteenId: "mensa_garching",
		Url:      "http://graphite-kom.srv.lrz.de/render/?from=-1hours&target=ap.ap*-?mg*.ssid.*&format=json",
		MaxCount: 1000}}
)

func (c *CronService) canteenHomometerCron() error {
	log.Info("Updating canteen homometer stats...")
	for _, canteen := range CANTEENS {
		log.Debug("Updating canteen homometer stats for: ", canteen.CanteenId)
		aps := requestApData(&canteen)
		if len(aps) <= 0 {
			log.Debug("No canteen homometer data points found for: ", canteen.CanteenId)
			continue
		}

		count := sumApCounts(aps)
		updateDb(&canteen, count, c.db)
		log.Debug("Canteen homometer stats updated for: ", canteen.CanteenId)
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
	resp, err := http.Get(canteen.Url)
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
