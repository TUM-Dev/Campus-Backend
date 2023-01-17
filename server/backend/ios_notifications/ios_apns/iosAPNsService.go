package ios_apns

import (
	"errors"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_logging"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	ErrCouldNotCreateTokenRequest = errors.New("could not create token request")
)

type Service struct {
	Repository *Repository
	Logger     *ios_logging.Service
}

func (s *Service) RequestGradeUpdateForDevice(deviceID string) error {
	campusRequestToken, err := s.Repository.CreateCampusTokenRequest(deviceID)

	if err != nil {
		log.Errorf("Could not create campus token request: %s", err)
		return ErrCouldNotCreateTokenRequest
	}

	notification := model.NewIOSNotificationPayload(deviceID).Background(campusRequestToken.RequestID, model.IOSBackgroundCampusTokenRequest)

	_, err = s.Repository.SendBackgroundNotification(notification)

	if err != nil {
		log.Errorf("Could not send background notification: %s", err)
		return ErrCouldNotSendNotification
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
