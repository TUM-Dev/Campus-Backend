package model

import (
	"fmt"
	"time"
)

const (
	IOSUpdateTypeGrades      = "grades"
	IOSMinimumUpdateInterval = 30
)

type IOSScheduledUpdateLog struct {
	ID        uint32    `gorm:"primary_key;auto_increment;not_null" json:"id"`
	DeviceID  string    `gorm:"index:idx_scheduled_update_log_device,unique" json:"deviceId"`
	Type      string    `gorm:"type:enum ('grades');" json:"type"`
	CreatedAt time.Time `gorm:"index:idx_scheduled_update_log_created,unique;autoCreateTime" json:"createdAt"`
}

func (log *IOSScheduledUpdateLog) IsGrades() bool {
	return log.Type == IOSUpdateTypeGrades
}

func (log *IOSScheduledUpdateLog) String() string {
	return fmt.Sprintf("IOSScheduledUpdateLog{ID: %d, DeviceID: %s, Type: %s, CreatedAt: %s}", log.ID, log.DeviceID, log.Type, log.CreatedAt)
}
