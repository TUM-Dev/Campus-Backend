package model

import "time"

// PushDeliveryLog deduplicates sent push notifications per device and logical event key.
type PushDeliveryLog struct {
	ID       int64     `gorm:"primaryKey;autoIncrement;column:id;type:bigint;not null;" json:"id"`
	DeviceID string    `gorm:"column:device_id;type:varchar(191);not null;uniqueIndex:idx_push_delivery_dedup,priority:1" json:"device_id"`
	DedupKey string    `gorm:"column:dedup_key;type:varchar(191);not null;uniqueIndex:idx_push_delivery_dedup,priority:2" json:"dedup_key"`
	SentAt   time.Time `gorm:"column:sent_at;type:timestamp;not null;default:current_timestamp();" json:"sent_at"`
}
