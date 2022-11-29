package backend

import (
	"context"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	ios "github.com/TUM-Dev/Campus-Backend/backend/ios_notifications"
	"github.com/TUM-Dev/Campus-Backend/backend/ios_notifications/ios_apns"
	"github.com/TUM-Dev/Campus-Backend/backend/ios_notifications/ios_scheduling"
	"github.com/TUM-Dev/Campus-Backend/backend/ios_notifications/ios_usage"
)

func (s *CampusServer) GetIOSNotificationsService() *ios.Service {
	repository := ios.NewRepository(s.db)

	return ios.NewService(repository)
}

func (s *CampusServer) GetIOSUsageService() *ios_usage.Service {
	repository := ios_usage.NewRepository(s.db)

	return ios_usage.NewService(repository)
}

func (s *CampusServer) GetIOSAPNsService() *ios_apns.Service {
	repository := ios_apns.NewRepository()

	return ios_apns.NewService(repository)
}

func (s *CampusServer) GetIOSSchedulingService() *ios_scheduling.Service {
	repository := ios_scheduling.NewRepository(s.db)

	return ios_scheduling.NewService(repository)
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

func (s *CampusServer) SendIOSTestNotification(ctx context.Context, req *pb.SendIOSTestNotificationRequest) (*pb.SendIOSTestNotificationReply, error) {
	service := s.GetIOSAPNsService()
	return service.SendTestNotification(req)
}
