package ios_apns

import (
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_logging"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Service struct {
	Repository *Repository
	Logger     *ios_logging.Service
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
		s.RequestGradeUpdateForDevice(device.DeviceID)
	}

	return &pb.SendIOSTestBackgroundNotificationReply{
		Message: "Test notifications sent",
	}, nil
}

func (s *Service) SendSuccessUpdateTestNotification(deviceID string) error {
	notification := model.NewIOSNotificationPayload(deviceID).Alert("Campus App Test Notification", "", "Getting grades is now easier than ever!")

	_, err := s.Repository.SendAlertNotification(notification)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) RequestGradeUpdateForDevice(deviceID string) error {
	campusRequestToken, err := s.Repository.CreateCampusTokenRequest(deviceID)

	if err != nil {
		log.Errorf("Could not create campus token request: %s", err)
		return err
	}

	notification := model.NewIOSNotificationPayload(deviceID).Background(campusRequestToken.RequestID, model.IOSBackgroundCampusTokenRequest)

	_, err = s.Repository.SendBackgroundNotification(notification)

	if err != nil {
		log.Errorf("Could not send background notification: %s", err)
		return err
	}

	s.Logger.LogTokenRequest("Token Requested: %s", campusRequestToken.RequestID)

	return nil
}

func NewService(repository *Repository) *Service {
	return &Service{
		Repository: repository,
		Logger:     ios_logging.NewLogger(&repository.DB),
	}
}

func NewCronService(db *gorm.DB) *Service {
	return NewService(NewCronRepository(db))
}
