package model

// NewsCacheVersion is a singleton incremented transactionally with feed edits.
type NewsCacheVersion struct {
	ID      uint8  `gorm:"column:id;type:tinyint unsigned;primaryKey;autoIncrement:false"`
	Version uint64 `gorm:"column:version;type:bigint unsigned;not null;default:0"`
}
