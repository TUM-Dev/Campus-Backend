package backend

import (
	"errors"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
)

func ValidateCreateDevice(request *pb.CreateDeviceRequest) error {
	if request.GetDeviceId() == "" {
		return errors.New("deviceId is empty")
	}

	if request.GetDeviceType() == pb.DeviceType_IOS && request.GetPublicKey() == "" {
		return errors.New("publicKey is needed for IOS")
	}

	return nil
}

func ValidateDeleteDevice(request *pb.DeleteDeviceRequest) error {
	if request.GetDeviceId() == "" {
		return errors.New("deviceId is empty")
	}

	return nil
}
