// Package apns contains the logic for sending push notifications to iOS devices.
// and communicating with the Apple Push Notification Service (APNs).
package apns

import (
	"errors"
	"os"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	ErrCouldNotCreateTokenRequest = errors.New("could not create token request")
)

type Service struct {
	Repository *Repository
	IsActive   bool
}

// RequestGradeUpdateForDevice stores a Request ID to the database and sends a background
// notification to the device with the given deviceID.
// The device will then send an update request to the server including the CampusToken
// and the request ID.
func (s *Service) RequestGradeUpdateForDevice(deviceID string) error {
	campusRequestToken, err := s.Repository.CreateCampusTokenRequest(deviceID)
	if err != nil {
		log.WithError(err).Error("Could not create campus token request")
		return ErrCouldNotCreateTokenRequest
	}

	notification := model.NewIOSNotificationPayload(deviceID).Background(campusRequestToken.RequestID, model.IOSBackgroundCampusTokenRequest)

	if _, err := s.Repository.SendNotification(notification, model.IOSAPNSPushTypeBackground); err != nil {
		log.WithError(err).Error("Could not send background notification")
		return errors.New("could not send notification")
	}
	return nil
}

func ValidateRequirementsForIOSNotificationsService() error {
	if os.Getenv("APNS_KEY_ID") == "" {
		return errors.New("APNS_KEY_ID env variable is not set")
	}

	if os.Getenv("APNS_TEAM_ID") == "" {
		return errors.New("APNS_TEAM_ID env variable is not set")
	}

	if os.Getenv("APNS_P8_FILE_PATH") == "" {
		return errors.New("APNS_P8_FILE_PATH env variable is not set")
	}

	if _, err := APNsEncryptionKeyFromFile(); err != nil {
		return errors.New("APNS P8 token is not valid or not set")
	}

	return nil
}

func NewService(repository *Repository) *Service {
	return &Service{
		Repository: repository,
		IsActive:   true,
	}
}

func NewCronService(db *gorm.DB) *Service {
	if repo, err := NewCronRepository(db); err != nil {
		return &Service{IsActive: false}
	} else {
		return NewService(repo)
	}
}
