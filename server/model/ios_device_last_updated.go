package model

import (
	"fmt"

	"github.com/guregu/null"
)

// IOSDeviceLastUpdated used as a result of a query that joins
// IOSDevice and IOSDeviceRequestLog tables.
type IOSDeviceLastUpdated struct {
	DeviceID    string    `json:"deviceId"`
	LastUpdated null.Time `json:"lastUpdated"`
	PublicKey   string    `json:"publicKey"`
}

func (device *IOSDeviceLastUpdated) String() string {
	time := "null"

	if device.LastUpdated.Valid {
		time = device.LastUpdated.Time.String()
	}
	return fmt.Sprintf("IOSDeviceLastUpdated{DeviceID: %s, LastUpdated: %s}", device.DeviceID, time)
}
