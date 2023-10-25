package backend

import (
	"errors"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
)

func ValidateCreateDevice(req *pb.CreateDeviceRequest) error {
	if req.DeviceId == "" {
		return errors.New("deviceId is empty")
	}

	if req.DeviceType == pb.DeviceType_IOS && req.PublicKey == "" {
		return errors.New("publicKey is needed for IOS")
	}

	return nil
}

func ValidateDeleteDevice(req *pb.DeleteDeviceRequest) error {
	if req.DeviceId == "" {
		return errors.New("deviceId is empty")
	}

	return nil
}
