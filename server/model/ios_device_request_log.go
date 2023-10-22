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
	DeviceID    string    `gorm:"size:200;not null" json:"deviceId"`
	Device      IOSDevice `gorm:"constraint:OnDelete:CASCADE;" json:"device"`
	RequestType string    `gorm:"not null;type:enum ('CAMPUS_TOKEN_REQUEST');" json:"requestType"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
}
