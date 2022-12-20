package model

import "time"

type IOSDevice struct {
	DeviceID  string    `gorm:"primary_key" json:"deviceId"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	PublicKey string    `gorm:"not null" json:"publicKey"`
}
