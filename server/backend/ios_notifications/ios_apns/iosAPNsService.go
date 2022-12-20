package ios_apns

import (
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"github.com/TUM-Dev/Campus-Backend/backend/ios_notifications"
	"github.com/TUM-Dev/Campus-Backend/backend/ios_notifications/ios_device"
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Repository *Repository
}

func (s *Service) SendTestNotification(request *pb.SendIOSTestNotificationRequest) (*pb.SendIOSTestNotificationReply, error) {
	if err := ios_notifications.ValidateSendTestNotification(request); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	devicesRepo := ios_device.NewRepository(&s.Repository.DB)

	devices, err := devicesRepo.GetDevices()

	if err != nil {
		return nil, status.Error(codes.Internal, "Could not get devices")
	}

	for _, device := range devices {
		notification := model.NewIOSNotificationPayload(device.DeviceID).Alert("Campus App Test Notification", "", "Getting grades is now easier than ever!")

		notification.Encrypt(device.PublicKey)

		_, err := s.Repository.SendAlertNotification(notification)

		if err != nil {
			return nil, err
		}
	}

	return &pb.SendIOSTestNotificationReply{
		Message: "Test notifications sent",
	}, nil
}

func (s *Service) SendTestBackgroundNotification() (*pb.SendIOSTestBackgroundNotificationReply, error) {
	devicesRepo := ios_device.NewRepository(&s.Repository.DB)

	devices, err := devicesRepo.GetDevices()

	if err != nil {
		return nil, status.Error(codes.Internal, "Could not get devices")
	}

	for _, device := range devices {
		campusRequestToken, err := s.Repository.CreateCampusTokenRequest(device.DeviceID)

		if err != nil {
			log.Errorf("Could not create campus token request: %s", err)
			return nil, err
		}

		notification := model.NewIOSNotificationPayload(device.DeviceID).Background(campusRequestToken.RequestID, model.IOSBackgroundCampusTokenRequest)

		_, err = s.Repository.SendBackgroundNotification(notification)

		if err != nil {
			return nil, err
		}
	}

	return &pb.SendIOSTestBackgroundNotificationReply{
		Message: "Test notifications sent",
	}, nil
}

func NewService(repository *Repository) *Service {
	return &Service{
		Repository: repository,
	}
}
