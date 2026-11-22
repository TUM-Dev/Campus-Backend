package cron

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// LibraryApInformation configures a TUM branch library tracked via LRZ Graphite Wi-Fi AP counts.
type LibraryApInformation struct {
	LibraryID string
	Name      string
	Campus    string
	Target    string
	MaxCount  uint32
}

// TrackedLibraries defines the TUM branch libraries monitored for live occupancy and busyness curves.
var TrackedLibraries = []LibraryApInformation{
	{
		LibraryID: "tb-stammgelaende",
		Name:      "Teilbibliothek Stammgelände",
		Campus:    "stammgelaende",
		Target:    "ap.ap*-?b1*.ssid.*",
		MaxCount:  650,
	},
	{
		LibraryID: "tb-maschinenwesen",
		Name:      "Teilbibliothek Maschinenwesen",
		Campus:    "garching",
		Target:    "ap.ap*-?mw*.ssid.*",
		MaxCount:  500,
	},
	{
		LibraryID: "tb-mathematik-informatik",
		Name:      "Teilbibliothek Mathematik & Informatik",
		Campus:    "garching",
		Target:    "ap.ap*-?mi*.ssid.*",
		MaxCount:  450,
	},
	{
		LibraryID: "tb-chemie",
		Name:      "Teilbibliothek Chemie",
		Campus:    "garching",
		Target:    "ap.ap*-?ch*.ssid.*",
		MaxCount:  250,
	},
	{
		LibraryID: "tb-physik",
		Name:      "Teilbibliothek Physik",
		Campus:    "garching",
		Target:    "ap.ap*-?ph*.ssid.*",
		MaxCount:  200,
	},
	{
		LibraryID: "tb-weihenstephan",
		Name:      "Teilbibliothek Weihenstephan",
		Campus:    "weihenstephan",
		Target:    "ap.ap*-?wb*.ssid.*",
		MaxCount:  300,
	},
	{
		LibraryID: "tb-sport-gesundheit",
		Name:      "Teilbibliothek Sport- und Gesundheitswissenschaften",
		Campus:    "olympiapark",
		Target:    "ap.ap*-?sg*.ssid.*",
		MaxCount:  180,
	},
}

func (c *CronService) libraryOccupancyCron() error {
	log.Debug("Updating library occupancy stats...")
	now := time.Now()
	for _, lib := range TrackedLibraries {
		if lib.Target == "" {
			continue
		}
		aps := requestLibraryApData(lib)
		if len(aps) == 0 {
			continue
		}
		count := sumApCounts(aps)
		if err := updateLibraryOccupancy(c.db, &lib, count, now); err != nil {
			log.WithError(err).WithField("library_id", lib.LibraryID).Warn("Failed to update library occupancy")
			continue
		}
		percent := clampedOccupancyPercent(count, lib.MaxCount)
		if err := recordLibraryOccupancyHistory(c.db, lib.LibraryID, count, percent, now); err != nil {
			log.WithError(err).WithField("library_id", lib.LibraryID).Warn("Failed to record library occupancy history")
		}
	}
	if err := pruneOldLibraryOccupancyHistory(c.db, now); err != nil {
		log.WithError(err).Warn("Failed to prune old library occupancy history")
	}
	return nil
}

func requestLibraryApData(lib LibraryApInformation) []AccessPoint {
	url := fmt.Sprintf(BaseUrl, lib.Target)
	resp, err := webfetch.Get(url)
	if err != nil {
		log.WithError(err).WithField("library_id", lib.LibraryID).Debug("Library occupancy web request failed")
		return nil
	}
	if resp.Body != nil {
		defer func(body io.ReadCloser) { _ = body.Close() }(resp.Body)
	}
	var aps []AccessPoint
	if err := json.NewDecoder(resp.Body).Decode(&aps); err != nil {
		log.WithError(err).WithField("library_id", lib.LibraryID).Warn("Library occupancy JSON decode failed")
		return nil
	}
	return aps
}

func updateLibraryOccupancy(db *gorm.DB, lib *LibraryApInformation, count uint32, now time.Time) error {
	entry := model.LibraryOccupancy{
		LibraryID: lib.LibraryID,
		Name:      lib.Name,
		Campus:    lib.Campus,
		Count:     count,
		MaxCount:  lib.MaxCount,
		Percent:   clampedOccupancyPercent(count, lib.MaxCount),
		Timestamp: now,
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "library_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "campus", "count", "max_count", "percent", "timestamp"}),
	}).Create(&entry).Error
}

func recordLibraryOccupancyHistory(db *gorm.DB, libraryID string, count uint32, percent float32, now time.Time) error {
	sampleDate, weekday, bucket := berlinOccupancyBucket(now)
	row := model.LibraryOccupancyHistory{
		LibraryID:   libraryID,
		SampleDate:  sampleDate,
		Bucket:      bucket,
		Weekday:     weekday,
		AvgCount:    float64(count),
		AvgPercent:  float64(percent),
		SampleCount: 1,
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "library_id"}, {Name: "sample_date"}, {Name: "bucket"}},
		DoUpdates: clause.Assignments(map[string]any{
			"avg_count":    gorm.Expr("(avg_count * sample_count + VALUES(avg_count)) / (sample_count + 1)"),
			"avg_percent":  gorm.Expr("(avg_percent * sample_count + VALUES(avg_percent)) / (sample_count + 1)"),
			"sample_count": gorm.Expr("sample_count + 1"),
		}),
	}).Create(&row).Error
}

func pruneOldLibraryOccupancyHistory(db *gorm.DB, now time.Time) error {
	cutoff, _, _ := berlinOccupancyBucket(now.AddDate(0, 0, -headCountHistoryRetentionDays))
	return db.Where("sample_date < ?", cutoff).Delete(&model.LibraryOccupancyHistory{}).Error
}
