package backend

import (
	"context"
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns/ios_apns_jwt"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_new_exams_callback"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_request_response"
	"google.golang.org/protobuf/types/known/emptypb"
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

	return ios_request_response.NewService(repository)
}

func (s *CampusServer) GetIOSNewExamsCallbackService() *ios_new_exams_callback.Service {
	repository := ios_new_exams_callback.NewRepository(s.db)

	return ios_new_exams_callback.NewService(repository, s.GetIOSAPNsService(), s.GetIOSNotificationsService().IsActive)
}

func (s *CampusServer) IOSDeviceRequestResponse(_ context.Context, req *pb.IOSDeviceRequestResponseRequest) (*pb.IOSDeviceRequestResponseReply, error) {
	service := s.GetIOSRequestResponseService()
	return service.HandleDeviceRequestResponse(req, s.iOSNotificationsService.IsActive)
}

func (s *CampusServer) IOSNewExamsHookCallback(_ context.Context, req *pb.NewExamsHookRequest) (*emptypb.Empty, error) {
	service := s.GetIOSNewExamsCallbackService()
	err := service.HandleNewExamsCallback(req)

	return &emptypb.Empty{}, err
}
