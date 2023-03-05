package model

import (
	"fmt"
	"time"
)

// IOSDevice stores relevant device information.
// E.g. the PublicKey which is used to encrypt push notifications
// The DeviceID can be used to send push notifications via APNs
type IOSDevice struct {
	DeviceID          string    `gorm:"primary_key" json:"deviceId"`
	CreatedAt         time.Time `gorm:"default:now()" json:"createdAt"`
	PublicKey         string    `gorm:"not null" json:"publicKey"`
	ActivityToday     int32     `gorm:"default:0" json:"activityToday"`
	ActivityThisWeek  int32     `gorm:"default:0" json:"activityThisWeek"`
	ActivityThisMonth int32     `gorm:"default:0" json:"activityThisMonth"`
	ActivityThisYear  int32     `gorm:"default:0" json:"activityThisYear"`
}

func (device *IOSDevice) String() string {
	return fmt.Sprintf("IOSDevice{DeviceID: %s}", device.DeviceID)
}
