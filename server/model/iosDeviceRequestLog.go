package model

import "time"

type IOSDeviceRequestLog struct {
	RequestID   string    `gorm:"primary_key;default:UUID()" json:"requestId"`
	DeviceID    string    `json:"deviceId" gorm:"size:200;not null"`
	Device      IOSDevice `json:"device"`
	RequestType string    `json:"requestType" gorm:"not null;type:enum ('CAMPUS_TOKEN_REQUEST');"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
