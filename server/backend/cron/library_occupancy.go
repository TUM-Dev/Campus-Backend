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
// AP targets and MaxCount capacities are derived from https://wlan.lrz.de/apstat/ reading-room APs.
// Note: Teilbibliothek Mathematik & Informatik (RBG Wi-Fi) and Medizin (MRI) are not reported in LRZ Graphite.
var TrackedLibraries = []LibraryApInformation{
	{
		// Geb. 0501 (Arcisstr. 21): Lesesäle 1136, 1140, 1142A, 1160, Z1136A, 2140, 2148, 2170
		LibraryID: "tb-stammgelaende",
		Name:      "Teilbibliothek Stammgelände",
		Campus:    "stammgelaende",
		Target:    "ap.{apa12-1bb,apa13-1bb,apa14-1bb,apa15-1bb,apa16-1bb,apa26-2bb,apa27-2bb,apa29-2bb,apa32-2bb,apa33-2bb,apa36-2bb,apa37-2bb,apa38-2bb,apa39-2bb,apa40-2bb,apa41-2bb}.ssid.*",
		MaxCount:  1200,
	},
	{
		// Geb. 5500 (Boltzmannstr. 15): Bauteil 0, 1025H Lesesaal/Monographien
		LibraryID: "tb-maschinenwesen",
		Name:      "Teilbibliothek Maschinenwesen",
		Campus:    "garching",
		Target:    "ap.apa{10,11,12,13}-1w0.ssid.*",
		MaxCount:  220,
	},
	{
		// Geb. 5401 (Lichtenbergstr. 4): 11140, 11150 Lesesaal, 21003 Scannerraum
		LibraryID: "tb-chemie",
		Name:      "Teilbibliothek Chemie",
		Campus:    "garching",
		Target:    "ap.{apa04-1af,apa09-1af,apa10-1af,apa13-2af}.ssid.*",
		MaxCount:  250,
	},
	{
		// Geb. 5101 (James-Franck-Str. 1): 3235, 4235 Lesesaal
		LibraryID: "tb-physik",
		Name:      "Teilbibliothek Physik",
		Campus:    "garching",
		Target:    "ap.{apa32-1at,apa33-1at,apa34-1at,apa41-2at,apa43-2at}.ssid.*",
		MaxCount:  160,
	},
	{
		// Geb. 4220 (Maximus-von-Imhof-Forum 3): Lesesaal EG/OG, Scannerraum, Gruppenarbeitsraum
		LibraryID: "tb-weihenstephan",
		Name:      "Teilbibliothek Weihenstephan",
		Campus:    "weihenstephan",
		Target:    "ap.{apa06-0qy,apa07-1qy,apa10-1qy,apa11-0qy,apa13-1qy}.ssid.*",
		MaxCount:  260,
	},
	{
		// Geb. 2927 (Petersgasse 18): Bibliothek Lesesaal 1 & 2, Garderobe
		LibraryID: "tb-straubing",
		Name:      "Teilbibliothek Straubing",
		Campus:    "straubing",
		Target:    "ap.apa{12,13,14}-0uc.ssid.*",
		MaxCount:  40,
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
