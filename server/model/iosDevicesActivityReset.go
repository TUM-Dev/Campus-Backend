package model

import "time"

const (
	IOSDevicesActivityResetTypeDay   = "day"
	IOSDevicesActivityResetTypeWeek  = "week"
	IOSDevicesActivityResetTypeMonth = "month"
	IOSDevicesActivityResetTypeYear  = "year"
)

type IOSDevicesActivityReset struct {
	Type      string    `gorm:"primary_key;type:enum('day', 'week', 'month', 'year')"`
	LastReset time.Time `gorm:"default:now()"`
}
