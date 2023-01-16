package model

import (
	"fmt"
	"time"
)

type IOSDevice struct {
	DeviceID          string    `gorm:"primary_key" json:"deviceId"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"createdAt"`
	PublicKey         string    `gorm:"not null" json:"publicKey"`
	activityToday     int32     `json:"activityToday"`
	activityThisWeek  int32     `json:"activityThisWeek"`
	activityThisMonth int32     `json:"activityThisMonth"`
	activityThisYear  int32     `json:"activityThisYear"`
}

func (device *IOSDevice) String() string {
	return fmt.Sprintf("IOSDevice{DeviceID: %s}", device.DeviceID)
}
