package backend

import (
	"context"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"github.com/TUM-Dev/Campus-Backend/backend/ios_notifications/ios_apns"
	"github.com/TUM-Dev/Campus-Backend/backend/ios_notifications/ios_apns/ios_apns_jwt"
	"github.com/TUM-Dev/Campus-Backend/backend/ios_notifications/ios_device"
	"github.com/TUM-Dev/Campus-Backend/backend/ios_notifications/ios_request_response"
	"github.com/TUM-Dev/Campus-Backend/backend/ios_notifications/ios_usage"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type IOSNotificationsService struct {
	DB        *gorm.DB
	APNSToken *ios_apns_jwt.Token
}

func (s *CampusServer) GetIOSDeviceService() *ios_device.Service {
	repository := ios_device.NewRepository(s.db)

	return ios_device.NewService(repository)
}

func (s *CampusServer) GetIOSUsageService() *ios_usage.Service {
	repository := ios_usage.NewRepository(s.db)

	return ios_usage.NewService(repository)
}

func (s *CampusServer) GetIOSAPNsService() *ios_apns.Service {
	repository := ios_apns.NewRepository(s.db, s.GetIOSNotificationsService().APNSToken)

	return ios_apns.NewService(repository)
}

func (s *CampusServer) GetIOSRequestResponseService() *ios_request_response.Service {
	repository := ios_request_response.NewRepository(s.db, s.GetIOSNotificationsService().APNSToken)

	return ios_request_response.NewService(repository)
}

func (s *CampusServer) RegisterIOSDevice(ctx context.Context, req *pb.RegisterIOSDeviceRequest) (*pb.RegisterIOSDeviceReply, error) {
	service := s.GetIOSDeviceService()
	return service.RegisterDevice(req)
}

func (s *CampusServer) RemoveIOSDevice(ctx context.Context, req *pb.RemoveIOSDeviceRequest) (*pb.RemoveIOSDeviceReply, error) {
	service := s.GetIOSDeviceService()
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

func (s *CampusServer) SendIOSTestBackgroundNotification(ctx context.Context, req *emptypb.Empty) (*pb.SendIOSTestBackgroundNotificationReply, error) {
	service := s.GetIOSAPNsService()
	return service.SendTestBackgroundNotification()
}

func (s *CampusServer) IOSDeviceRequestResponse(ctx context.Context, req *pb.IOSDeviceRequestResponseRequest) (*pb.IOSDeviceRequestResponseReply, error) {
	service := s.GetIOSRequestResponseService()
	return service.HandleDeviceRequestResponse(req)
}
