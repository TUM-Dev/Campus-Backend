package model

import "time"

// IOSDeviceUsageLog is the model for the ios_device_usage_log table.
// Probably **DEPRECATED**, will be replaced by the activity parameters
// in the IOSDevice model.
type IOSDeviceUsageLog struct {
	ID        uint32    `gorm:"primary_key;auto_increment;not_null" json:"id"`
	DeviceID  string    `gorm:"index:idx_usage_log_created,unique" json:"deviceId"`
	CreatedAt time.Time `gorm:"index:idx_usage_log_created,unique;autoCreateTime" json:"createdAt"`
}
