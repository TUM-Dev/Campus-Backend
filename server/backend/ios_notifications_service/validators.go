package ios_notifications_service

import (
	"errors"
	pb "github.com/TUM-Dev/Campus-Backend/api"
)

func ValidateRegisterDevice(request *pb.RegisterIOSDeviceRequest) error {
	if request.GetDeviceId() == "" {
		return errors.New("deviceId is empty")
	}
	return nil
}

func ValidateRemoveDevice(request *pb.RemoveIOSDeviceRequest) error {
	if request.GetDeviceId() == "" {
		return errors.New("deviceId is empty")
	}

	return nil
}
