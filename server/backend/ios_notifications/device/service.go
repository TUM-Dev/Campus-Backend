// Package device provides functions to create/delete ios devices
package device

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
	iosRegisteredDevices = promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "ios",
		Name:      "ios_created_devices",
		Help:      "The number of currently created ios devices",
	})
)

func (service *Service) CreateDevice(req *pb.CreateDeviceRequest) (*pb.CreateDeviceReply, error) {
	device := model.IOSDevice{
		DeviceID:  req.DeviceId,
		PublicKey: req.PublicKey,
	}

	if err := service.Repository.CreateDevice(&device); err != nil {
		return nil, status.Error(codes.Internal, "Could not create device")
	}
	iosRegisteredDevices.Inc()

	return &pb.CreateDeviceReply{
		DeviceId: device.DeviceID,
	}, nil
}

func (service *Service) DeleteDevice(req *pb.DeleteDeviceRequest) (*pb.DeleteDeviceReply, error) {
	if err := service.Repository.DeleteDevice(req.DeviceId); err != nil {
		return nil, status.Error(codes.Internal, "Could not delete device")
	}

	iosRegisteredDevices.Dec()
	return &pb.DeleteDeviceReply{
		DeviceId: req.DeviceId,
	}, nil
}

func NewService(db *Repository) *Service {
	return &Service{
		Repository: db,
	}
}
