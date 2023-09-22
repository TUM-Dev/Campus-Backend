// Package ios_device provides functions to register and remove ios devices
package ios_device

import (
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Service struct {
	Repository *Repository
}

var (
	ErrCouldNotRegisterDevice = status.Error(codes.Internal, "Could not register device")
	ErrCouldNotRemoveDevice   = status.Error(codes.Internal, "Could not remove device")

	iosRegisteredDevices = promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "ios",
		Name:      "ios_registered_devices",
		Help:      "The number of currently registered ios devices",
	})
)

func (service *Service) RegisterDevice(request *pb.RegisterDeviceRequest) (*pb.RegisterDeviceReply, error) {
	device := model.IOSDevice{
		DeviceID:  request.GetDeviceId(),
		PublicKey: request.GetPublicKey(),
	}

	if err := service.Repository.RegisterDevice(&device); err != nil {
		return nil, ErrCouldNotRegisterDevice
	}
	iosRegisteredDevices.Inc()

	return &pb.RegisterDeviceReply{
		DeviceId: device.DeviceID,
	}, nil
}

func (service *Service) RemoveDevice(request *pb.RemoveDeviceRequest) (*pb.RemoveDeviceReply, error) {
	if err := service.Repository.RemoveDevice(request.GetDeviceId()); err != nil {
		return nil, ErrCouldNotRemoveDevice
	}

	iosRegisteredDevices.Dec()
	return &pb.RemoveDeviceReply{
		DeviceId: request.GetDeviceId(),
	}, nil
}

func NewService(db *Repository) *Service {
	return &Service{
		Repository: db,
	}
}
