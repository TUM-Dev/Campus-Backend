package model

import "time"

// LibraryOccupancy stores the latest live Wi-Fi AP occupancy estimate for a TUM library.
type LibraryOccupancy struct {
	LibraryID string    `gorm:"primaryKey;column:library_id;type:varchar(64);not null;" json:"library_id"`
	Name      string    `gorm:"column:name;type:varchar(128);not null;" json:"name"`
	Campus    string    `gorm:"column:campus;type:varchar(64);not null;" json:"campus"`
	Count     uint32    `gorm:"column:count;type:int unsigned;not null;" json:"count"`
	MaxCount  uint32    `gorm:"column:max_count;type:int unsigned;not null;" json:"max_count"`
	Percent   float32   `gorm:"column:percent;type:float;not null;" json:"percent"`
	Timestamp time.Time `gorm:"column:timestamp;type:timestamp;not null;default:current_timestamp();" json:"timestamp"`
}

// LibraryOccupancyHistory stores 15-minute occupancy rollups per library and day over an 8-week window.
type LibraryOccupancyHistory struct {
	LibraryID   string  `gorm:"primaryKey;column:library_id;type:varchar(64);not null;index:idx_library_occ_history_lookup,priority:1" json:"library_id"`
	SampleDate  string  `gorm:"primaryKey;column:sample_date;type:char(10);not null;index:idx_library_occ_history_date" json:"sample_date"`
	Bucket      int32   `gorm:"primaryKey;column:bucket;type:tinyint;not null;index:idx_library_occ_history_lookup,priority:3" json:"bucket"`
	Weekday     int32   `gorm:"column:weekday;type:tinyint;not null;index:idx_library_occ_history_lookup,priority:2" json:"weekday"`
	AvgCount    float64 `gorm:"column:avg_count;type:double;not null;" json:"avg_count"`
	AvgPercent  float64 `gorm:"column:avg_percent;type:double;not null;" json:"avg_percent"`
	SampleCount int32   `gorm:"column:sample_count;type:int;not null;default:1;" json:"sample_count"`
}

func (LibraryOccupancyHistory) TableName() string {
	return "library_occupancy_history"
}
