package cron

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccessPoint struct {
	Target     string      `json:"target"`
	DataPoints [][]float32 `json:"datapoints"`
}

type CanteenApInformation struct {
	// The canteen ID defined by the ea-api
	CanteenId string
	// Graphite access point filter list
	Target string
	// An approximated amount of people that once present indicate the canteen is full/running out of capacity
	MaxCount uint32
}

var (
	Canteens = []CanteenApInformation{
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
		// No data found: http://wlan.lrz.de/apstat/filter/Unterbezirk/w0/
		{CanteenId: "stucafe-boltzmannstr",
			Target:   "",
			MaxCount: 1000},
		// No APs for the stucafe found: http://wlan.lrz.de/apstat/filter/Unterbezirk/ef/
		{CanteenId: "stucafe-connollystr",
			Target:   "",
			MaxCount: 1000},
		{CanteenId: "stucafe-garching",
			Target:   "ap.apa{14,15,16}-?mg*.ssid.*",
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
BaseUrl is the base URL for the required format.
Contains the '%s' placeholder that has to replaced with the Target property of the
CanteenApInformation when performing a request.
*/
const BaseUrl = "http://graphite-kom.srv.lrz.de/render/?from=-10min&target=%s&format=json"

const headCountHistoryRetentionDays = 180 // ~6 months

func (c *CronService) canteenHeadCountCron() error {
	log.Debug("Updating canteen head count stats...")
	now := time.Now()
	for _, canteen := range Canteens {
		if len(canteen.Target) <= 0 {
			log.WithField("CanteenId", canteen.CanteenId).Debug("Skipping canteen head count stats, since there is no target.")
			continue
		}

		log.WithField("CanteenId", canteen.CanteenId).Debug("Updating canteen head count stats")
		aps := canteen.requestApData()
		if len(aps) <= 0 {
			log.WithField("CanteenId", canteen.CanteenId).Debug("No canteen head count data points found")
			continue
		}

		count := sumApCounts(aps)
		fields := log.Fields{"count": count, "CanteenId": canteen.CanteenId}
		if err := updateDb(&canteen, count, c.db); err != nil {
			log.WithFields(fields).WithError(err).Error("Failed to update Canteen head count stats")
			continue
		}
		percent := clampedOccupancyPercent(count, canteen.MaxCount)
		if err := recordHeadCountHistory(c.db, canteen.CanteenId, count, percent, now); err != nil {
			log.WithFields(fields).WithError(err).Warn("Failed to record canteen head count history")
		}
		log.WithFields(fields).Debug("Canteen head count stats updated")
	}
	if err := pruneOldHeadCountHistory(c.db, now); err != nil {
		log.WithError(err).Warn("Failed to prune old canteen head count history")
	}
	log.Debug("Canteen head count stats updated.")
	return nil
}

func clampedOccupancyPercent(count, maxCount uint32) float32 {
	if maxCount == 0 {
		return 0
	}
	percent := (float32(count) / float32(maxCount)) * 100
	if percent > 100.0 {
		return 100.0
	}
	return percent
}

func berlinOccupancyBucket(now time.Time) (sampleDate string, weekday int32, bucket int32) {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		loc = time.UTC
	}
	local := now.In(loc)
	sampleDate = local.Format("2006-01-02")
	weekday = int32((int(local.Weekday()) + 6) % 7) // 0=Monday .. 6=Sunday
	bucket = int32(local.Hour()*4 + local.Minute()/15)
	return sampleDate, weekday, bucket
}

func recordHeadCountHistory(db *gorm.DB, canteenID string, count uint32, percent float32, now time.Time) error {
	sampleDate, weekday, bucket := berlinOccupancyBucket(now)
	row := model.CanteenHeadCountHistory{
		CanteenID:   canteenID,
		SampleDate:  sampleDate,
		Bucket:      bucket,
		Weekday:     weekday,
		AvgCount:    float64(count),
		AvgPercent:  float64(percent),
		SampleCount: 1,
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "canteen_id"}, {Name: "sample_date"}, {Name: "bucket"}},
		DoUpdates: clause.Assignments(map[string]any{
			"avg_count":    gorm.Expr("(avg_count * sample_count + VALUES(avg_count)) / (sample_count + 1)"),
			"avg_percent":  gorm.Expr("(avg_percent * sample_count + VALUES(avg_percent)) / (sample_count + 1)"),
			"sample_count": gorm.Expr("sample_count + 1"),
		}),
	}).Create(&row).Error
}

func pruneOldHeadCountHistory(db *gorm.DB, now time.Time) error {
	cutoff, _, _ := berlinOccupancyBucket(now.AddDate(0, 0, -headCountHistoryRetentionDays))
	return db.Where("sample_date < ?", cutoff).Delete(&model.CanteenHeadCountHistory{}).Error
}

func sumApCounts(aps []AccessPoint) uint32 {
	total := uint32(0)
	for _, ap := range aps {
		_, count, err := parseAp(&ap)
		if err != nil {
			log.WithError(err).Error("Canteen HeadCount getting the count failed for access point: ", ap.Target)
			continue
		}
		total += count
	}
	return total
}

func updateDb(canteen *CanteenApInformation, count uint32, db *gorm.DB) error {
	percent := (float32(count) / float32(canteen.MaxCount)) * 100
	// Ensure we do not see more than 100%
	if percent > 100.0 {
		percent = 100.0
	}
	entry := model.CanteenHeadCount{
		CanteenId: canteen.CanteenId,
		Count:     count,
		MaxCount:  canteen.MaxCount,
		Percent:   percent,
		Timestamp: time.Now(),
	}

	res := db.Model(&model.CanteenHeadCount{}).Where(model.CanteenHeadCount{CanteenId: canteen.CanteenId}).
		Select("Count", "MaxCount", "Percent", "Timestamp").Updates(&entry)
	if res.Error != nil {
		log.WithError(res.Error).WithField("CanteenId", canteen.CanteenId).Error("could not update all instances of headcount")
		return res.Error
	}

	if res.RowsAffected == 0 {
		if err := db.Create(&entry).Error; err != nil {
			fields := log.Fields{
				"CanteenId": entry.CanteenId,
				"Count":     entry.Count,
				"MaxCount":  entry.MaxCount,
				"Percent":   entry.Percent,
				"Timestamp": entry.Timestamp}
			log.WithError(res.Error).WithFields(fields).Error("could not create headcount entry")
			return err
		}
	}
	return nil
}

func (canteen CanteenApInformation) requestApData() []AccessPoint {
	// Perform web request
	url := fmt.Sprintf(BaseUrl, canteen.Target)
	resp, err := webfetch.Get(url)
	if err != nil {
		log.WithError(err).WithField("CanteenId", canteen.CanteenId).Error("Canteen HeadCount web request failed")
		return []AccessPoint{}
	}

	// Ensure we close the body once we leave this function
	if resp.Body != nil {
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				log.WithError(err).Error("Could not close body")
			}
		}(resp.Body)
	}

	// Parse as JSON
	var aps []AccessPoint
	if err = json.NewDecoder(resp.Body).Decode(&aps); err != nil {
		log.WithError(err).WithField("CanteenId", canteen.CanteenId).Error("Canteen HeadCount parsing failed")
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
