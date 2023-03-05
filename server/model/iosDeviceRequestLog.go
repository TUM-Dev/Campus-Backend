package model

import "time"

// An IOSDeviceRequestLog is created when the backend wants to request data from the device.
//
// 1. The backend creates a new IOSDeviceRequestLog
//
// 2. The backend sends a background push notification to the device containing
// the RequestID of the IOSDeviceRequestLog.
//
// 3. The device receives the push notification and sends a request to the backend
// containing the RequestID and the data.
type IOSDeviceRequestLog struct {
	RequestID   string    `gorm:"primary_key;default:UUID()" json:"requestId"`
	DeviceID    string    `json:"deviceId" gorm:"size:200;not null"`
	Device      IOSDevice `json:"device" gorm:"constraint:OnDelete:CASCADE;"`
	RequestType string    `json:"requestType" gorm:"not null;type:enum ('CAMPUS_TOKEN_REQUEST');"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
