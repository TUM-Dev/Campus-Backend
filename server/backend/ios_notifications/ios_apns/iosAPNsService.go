package ios_apns

import (
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"github.com/TUM-Dev/Campus-Backend/model"
)

type Service struct {
	Repository *Repository
}

func (s *Service) SendTestNotification(request *pb.SendIOSTestNotificationRequest) (*pb.SendIOSTestNotificationReply, error) {
	notification := model.IOSRemoteAlertNotification{
		DeviceId: request.DeviceId,
		Aps: model.IOSAlertAPS{
			Alert: model.IOSAlertAPSContent{
				Title:    "Test Notification",
				Subtitle: request.Message,
				Body:     "This is a test notification body",
			},
		},
	}

	_, err := s.Repository.SendTestNotification(&notification)

	if err != nil {
		return nil, err
	}

	return &pb.SendIOSTestNotificationReply{
		Message: "Test notification sent",
	}, nil
}

func NewService(repository *Repository) *Service {
	return &Service{
		Repository: repository,
	}
}
