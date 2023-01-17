package ios_notifications

import (
	"errors"
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
)

func ValidateRegisterDevice(request *pb.RegisterIOSDeviceRequest) error {
	if request.GetDeviceId() == "" {
		return errors.New("deviceId is empty")
	}

	if request.GetPublicKey() == "" {
		return errors.New("publicKey is empty")
	}

	return nil
}

func ValidateRemoveDevice(request *pb.RemoveIOSDeviceRequest) error {
	if request.GetDeviceId() == "" {
		return errors.New("deviceId is empty")
	}

	return nil
}
