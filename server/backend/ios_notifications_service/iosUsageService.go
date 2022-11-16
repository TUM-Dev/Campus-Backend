package ios_notifications_service

import (
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"github.com/TUM-Dev/Campus-Backend/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IOSUsageService struct {
	Repository *IOSUsageRepository
}

func (service *IOSUsageService) AddUsage(request *pb.AddIOSDeviceUsageRequest) (*pb.AddIOSDeviceUsageReply, error) {
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
