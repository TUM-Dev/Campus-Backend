package ios_device

import (
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications"
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

func (service *Service) RegisterDevice(request *pb.RegisterIOSDeviceRequest) (*pb.RegisterIOSDeviceReply, error) {
	if err := ios_notifications.ValidateRegisterDevice(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	device := model.IOSDevice{
		DeviceID:  request.GetDeviceId(),
		PublicKey: request.GetPublicKey(),
	}

	err := service.Repository.RegisterDevice(&device)

	if err != nil {
		return nil, ErrCouldNotRegisterDevice
	}

	return &pb.RegisterIOSDeviceReply{
		DeviceId: device.DeviceID,
	}, nil
}

func (service *Service) RemoveDevice(request *pb.RemoveIOSDeviceRequest) (*pb.RemoveIOSDeviceReply, error) {

	if err := ios_notifications.ValidateRemoveDevice(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := service.Repository.RemoveDevice(request.GetDeviceId())

	if err != nil {
		return nil, ErrCouldNotRemoveDevice
	}

	return &pb.RemoveIOSDeviceReply{
		Message: "Successfully removed device",
	}, nil
}

func NewService(db *Repository) *Service {
	return &Service{
		Repository: db,
	}
}
