package ios_usage

import (
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Repository *Repository
}

func (service *Service) AddUsage(request *pb.AddIOSDeviceUsageRequest) (*pb.AddIOSDeviceUsageReply, error) {
	usageLog := model.IOSDeviceUsageLog{
		DeviceID: request.GetDeviceId(),
	}

	usage, err := service.Repository.AddUsage(&usageLog)

	if err != nil {
		return nil, status.Error(codes.Internal, "Could not add the device usage")
	}

	return &pb.AddIOSDeviceUsageReply{
		DeviceId:  usage.DeviceID,
		CreatedAt: usage.CreatedAt.String(),
		Id:        usage.ID,
	}, nil
}

func NewService(repository *Repository) *Service {
	return &Service{
		Repository: repository,
	}
}
