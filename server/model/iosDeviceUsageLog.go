package model

import "time"

type IOSDeviceUsageLog struct {
	ID        uint32    `gorm:"primary_key;auto_increment;not_null" json:"id"`
	DeviceID  string    `gorm:"index:idx_usage_log_created,unique" json:"deviceId"`
	CreatedAt time.Time `gorm:"index:idx_usage_log_created,unique;autoCreateTime" json:"createdAt"`
}
