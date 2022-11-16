package backend

import (
	"context"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	ios "github.com/TUM-Dev/Campus-Backend/backend/ios_notifications_service"
)

func (s *CampusServer) GetIOSNotificationsService() ios.IOSNotificationsService {
	repository := ios.IOSNotificationsRepository{
		DB: s.db,
	}

	return ios.IOSNotificationsService{Repository: &repository}
}

func (s *CampusServer) GetIOSUsageService() ios.IOSUsageService {
	repository := ios.IOSUsageRepository{
		DB: s.db,
	}

	return ios.IOSUsageService{Repository: &repository}
}

func (s *CampusServer) RegisterIOSDevice(ctx context.Context, req *pb.RegisterIOSDeviceRequest) (*pb.RegisterIOSDeviceReply, error) {
	service := s.GetIOSNotificationsService()
	return service.RegisterDevice(req)
}

func (s *CampusServer) RemoveIOSDevice(ctx context.Context, req *pb.RemoveIOSDeviceRequest) (*pb.RemoveIOSDeviceReply, error) {
	service := s.GetIOSNotificationsService()
	return service.RemoveDevice(req)
}

func (s *CampusServer) AddIOSDeviceUsage(ctx context.Context, req *pb.AddIOSDeviceUsageRequest) (*pb.AddIOSDeviceUsageReply, error) {
	service := s.GetIOSUsageService()
	return service.AddUsage(req)
}
