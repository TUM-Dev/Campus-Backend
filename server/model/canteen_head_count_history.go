package model

// CanteenHeadCountHistory stores aggregated 15-minute occupancy buckets per canteen and day
// to compute typical weekday busyness curves over a rolling 8-week window.
type CanteenHeadCountHistory struct {
	CanteenID   string  `gorm:"primaryKey;column:canteen_id;type:varchar(64);not null;index:idx_canteen_hc_history_lookup,priority:1" json:"canteen_id"`
	SampleDate  string  `gorm:"primaryKey;column:sample_date;type:char(10);not null;index:idx_canteen_hc_history_date" json:"sample_date"`
	Bucket      int32   `gorm:"primaryKey;column:bucket;type:tinyint;not null;index:idx_canteen_hc_history_lookup,priority:3" json:"bucket"`
	Weekday     int32   `gorm:"column:weekday;type:tinyint;not null;index:idx_canteen_hc_history_lookup,priority:2" json:"weekday"`
	AvgCount    float64 `gorm:"column:avg_count;type:double;not null;" json:"avg_count"`
	AvgPercent  float64 `gorm:"column:avg_percent;type:double;not null;" json:"avg_percent"`
	SampleCount int32   `gorm:"column:sample_count;type:int;not null;default:1;" json:"sample_count"`
}

func (CanteenHeadCountHistory) TableName() string {
	return "canteen_head_count_history"
}
