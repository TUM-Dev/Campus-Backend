package ios_notifications_service

import (
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"github.com/TUM-Dev/Campus-Backend/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IOSNotificationsService struct {
	Repository *IOSNotificationsRepository
}

func (service *IOSNotificationsService) RegisterDevice(request *pb.RegisterIOSDeviceRequest) (*pb.RegisterIOSDeviceReply, error) {

	if err := ValidateRegisterDevice(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	device := model.IOSDevice{
		DeviceID: request.GetDeviceId(),
	}

	err := service.Repository.RegisterDevice(&device)

	if err != nil {
		return nil, status.Error(codes.Internal, "Could not register device")
	}

	return &pb.RegisterIOSDeviceReply{
		DeviceId: device.DeviceID,
	}, nil
}

func (service *IOSNotificationsService) RemoveDevice(request *pb.RemoveIOSDeviceRequest) (*pb.RemoveIOSDeviceReply, error) {

	if err := ValidateRemoveDevice(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := service.Repository.RemoveDevice(request.GetDeviceId())

	if err != nil {
		return nil, status.Error(codes.Internal, "Could not remove device")
	}

	return &pb.RemoveIOSDeviceReply{
		Message: "Successfully removed device",
	}, nil
}

// func (service *IOSNotificationsService) GetDevices() (pb.)
