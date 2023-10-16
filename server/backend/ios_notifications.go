package backend

import (
	"context"
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/backend/http_header_authorization"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/apns"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_new_exams_callback"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/request_response"
	"gorm.io/gorm"
)

type IOSNotificationsService struct {
	DB        *gorm.DB
	APNSToken *apns.JWTToken
	IsActive  bool
}

func (s *CampusServer) GetIOSDeviceService() *device.Service {
	repository := device.NewRepository(s.db)

	return device.NewService(repository)
}

func (s *CampusServer) GetIOSAPNsService() *apns.Service {
	repository := apns.NewRepository(s.db, s.GetIOSNotificationsService().APNSToken)

	return apns.NewService(repository)
}

func (s *CampusServer) GetIOSRequestResponseService() *request_response.Service {
	repository := request_response.NewRepository(s.db, s.GetIOSNotificationsService().APNSToken)

	return request_response.NewService(repository)
}

func (s *CampusServer) GetIOSNewExamsCallbackService() *ios_new_exams_callback.Service {
	return ios_new_exams_callback.NewService(s.GetIOSAPNsService(), s.db, s.GetIOSNotificationsService().IsActive)
}

func (s *CampusServer) IOSDeviceRequestResponse(_ context.Context, req *pb.IOSDeviceRequestResponseRequest) (*pb.IOSDeviceRequestResponseReply, error) {
	service := s.GetIOSRequestResponseService()
	return service.HandleDeviceRequestResponse(req, s.iOSNotificationsService.IsActive)
}

func (s *CampusServer) IOSNewExamsHookCallback(ctx context.Context, req *pb.NewExamsHookRequest) (*pb.NewExamsHookReply, error) {
	err := http_header_authorization.CheckApiKeyAuthorization(ctx)
	if err != nil {
		return &pb.NewExamsHookReply{}, err
	}

	service := s.GetIOSNewExamsCallbackService()
	err = service.HandleNewExamsCallback(req)

	return &pb.NewExamsHookReply{}, err
}
