// Package ios_apns contains the logic for sending push notifications to iOS devices.
// and communicating with the Apple Push Notification Service (APNs).
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
}

// RequestGradeUpdateForDevice stores a Request ID to the database and sends a background
// notification to the device with the given deviceID.
// The device will then send an update request to the server including the CampusToken
// and the request ID.
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

	return nil
}

func NewService(repository *Repository) *Service {
	return &Service{
		Repository: repository,
	}
}

func NewCronService(db *gorm.DB) *Service {
	return NewService(NewCronRepository(db))
}
