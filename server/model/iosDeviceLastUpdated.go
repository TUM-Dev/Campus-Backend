package model

import (
	"database/sql"
	"fmt"
)

type IOSDeviceLastUpdated struct {
	DeviceID    string       `json:"deviceId"`
	LastUpdated sql.NullTime `json:"lastUpdated"`
	PublicKey   string       `json:"publicKey"`
}

func (device *IOSDeviceLastUpdated) String() string {
	return fmt.Sprintf("IOSDeviceLastUpdated{DeviceID: %s, LastUpdated: %d}", device.DeviceID, device.LastUpdated)
}
