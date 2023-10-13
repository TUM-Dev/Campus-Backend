package model

import (
	"fmt"
	"time"
)

// IOSDevice stores relevant device information.
// E.g. the PublicKey which is used to encrypt push notifications
// The DeviceID can be used to send push notifications via APNs
type IOSDevice struct {
	DeviceID  string    `gorm:"primary_key" json:"deviceId"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	PublicKey string    `gorm:"not null" json:"publicKey"`
}

func (device *IOSDevice) String() string {
	return fmt.Sprintf("IOSDevice{DeviceID: %s}", device.DeviceID)
}
