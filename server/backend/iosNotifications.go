package backend

import (
	"context"
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns/ios_apns_jwt"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_request_response"
	"gorm.io/gorm"
)

type IOSNotificationsService struct {
	DB        *gorm.DB
	APNSToken *ios_apns_jwt.Token
	IsActive  bool
}

func (s *CampusServer) GetIOSDeviceService() *ios_device.Service {
	repository := ios_device.NewRepository(s.db)

	return ios_device.NewService(repository)
}

func (s *CampusServer) GetIOSAPNsService() *ios_apns.Service {
	repository := ios_apns.NewRepository(s.db, s.GetIOSNotificationsService().APNSToken)

	return ios_apns.NewService(repository)
}

func (s *CampusServer) GetIOSRequestResponseService() *ios_request_response.Service {
	repository := ios_request_response.NewRepository(s.db, s.GetIOSNotificationsService().APNSToken)

	return ios_request_response.NewService(repository, s.GetIOSAPNsService())
}

func (s *CampusServer) IOSDeviceRequestResponse(_ context.Context, req *pb.IOSDeviceRequestResponseRequest) (*pb.IOSDeviceRequestResponseReply, error) {
	service := s.GetIOSRequestResponseService()
	return service.HandleDeviceRequestResponse(req, s.iOSNotificationsService.IsActive)
}
