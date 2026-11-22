package backend

import (
	"errors"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
)

func ValidateCreateDevice(req *pb.CreateDeviceRequest) error {
	if req.DeviceId == "" {
		return errors.New("deviceId is empty")
	}

	if req.FcmToken == "" && req.PublicKey == "" {
		return errors.New("fcm_token or public_key is required")
	}

	return nil
}

func ValidateDeleteDevice(req *pb.DeleteDeviceRequest) error {
	if req.DeviceId == "" {
		return errors.New("deviceId is empty")
	}

	return nil
}
