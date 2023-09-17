// Package ios_device provides functions to register and remove ios devices
package ios_device

import (
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/backend/influx"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Repository *Repository
}

var (
	ErrCouldNotRegisterDevice = status.Error(codes.Internal, "Could not register device")
	ErrCouldNotRemoveDevice   = status.Error(codes.Internal, "Could not remove device")
)

func (service *Service) RegisterDevice(request *pb.RegisterDeviceRequest) (*pb.RegisterDeviceReply, error) {
	device := model.IOSDevice{
		DeviceID:  request.GetDeviceId(),
		PublicKey: request.GetPublicKey(),
	}

	err := service.Repository.RegisterDevice(&device)

	if err != nil {
		return nil, ErrCouldNotRegisterDevice
	}

	influx.LogIOSRegisterDevice(request.GetDeviceId())

	return &pb.RegisterDeviceReply{
		DeviceId: device.DeviceID,
	}, nil
}

func (service *Service) RemoveDevice(request *pb.RemoveDeviceRequest) (*pb.RemoveDeviceReply, error) {
	err := service.Repository.RemoveDevice(request.GetDeviceId())

	if err != nil {
		return nil, ErrCouldNotRemoveDevice
	}

	influx.LogIOSRemoveDevice(request.GetDeviceId())

	return &pb.RemoveDeviceReply{
		DeviceId: request.GetDeviceId(),
	}, nil
}

func NewService(db *Repository) *Service {
	return &Service{
		Repository: db,
	}
}
