package model

import "time"

const (
	IOSUpdateTypeGrades = "grades"
)

type IOSScheduledUpdateLog struct {
	ID        uint32    `gorm:"primary_key;auto_increment;not_null" json:"id"`
	DeviceID  string    `gorm:"index:idx_scheduled_update_log_device,unique" json:"deviceId"`
	Type      string    `gorm:"type:enum ('grades');" json:"type"`
	CreatedAt time.Time `gorm:"index:idx_scheduled_update_log_created,unique;autoCreateTime" json:"createdAt"`
}
