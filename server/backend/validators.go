package backend

import (
	"errors"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
)

func ValidateRegisterDevice(request *pb.RegisterDeviceRequest) error {
	if request.GetDeviceId() == "" {
		return errors.New("deviceId is empty")
	}

	if request.GetDeviceType() == pb.DeviceType_IOS && request.GetPublicKey() == "" {
		return errors.New("publicKey is needed for IOS")
	}

	return nil
}

func ValidateRemoveDevice(request *pb.RemoveDeviceRequest) error {
	if request.GetDeviceId() == "" {
		return errors.New("deviceId is empty")
	}

	return nil
}
