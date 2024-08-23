package model

import (
	"github.com/guregu/null"
	"time"
)

// WifiMeasurement stores Wi-Fi measurements gathered by the lrz
type WifiMeasurement struct {
	Id               int64     `gorm:"primary_key;AUTO_INCREMENT"`
	Date             time.Time `gorm:"type:date;not null"`
	SSID             string    `gorm:"type:varchar(32);not null"`
	BSSID            string    `gorm:"type:varchar(64);not null"`
	DBm              null.Int
	AccuracyInMeters float32
	Latitude         float64
	Longitude        float64
}
