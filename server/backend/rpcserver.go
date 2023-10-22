package backend

import (
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/apns"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CampusServer struct {
	pb.UnimplementedCampusServiceServer
	pb.UnimplementedCanteenServiceServer
	pb.UnimplementedNewsServiceServer
	pb.UnimplementedLegacyServiceServer
	pb.UnimplementedNotificationServiceServer
	db                      *gorm.DB
	deviceBuf               *deviceBuffer // deviceBuf stores all devices from recent request and flushes them to db
	iOSNotificationsService *IOSNotificationsService
}

// Verify that CampusServer implements the required interfaces
var _ pb.CampusServiceServer = (*CampusServer)(nil)
var _ pb.CanteenServiceServer = (*CampusServer)(nil)
var _ pb.NewsServiceServer = (*CampusServer)(nil)
var _ pb.LegacyServiceServer = (*CampusServer)(nil)
var _ pb.NotificationServiceServer = (*CampusServer)(nil)

func New(db *gorm.DB) *CampusServer {
	log.Trace("Server starting up")
	return &CampusServer{
		db:                      db,
		deviceBuf:               newDeviceBuffer(),
		iOSNotificationsService: NewIOSNotificationsService(),
	}
}

func NewIOSNotificationsService() *IOSNotificationsService {
	if err := apns.ValidateRequirementsForIOSNotificationsService(); err != nil {
		log.WithError(err).Warn("failed to validate requirements for ios notifications service")

		return &IOSNotificationsService{
			APNSToken: nil,
			IsActive:  false,
		}
	}

	token, err := apns.NewToken()
	if err != nil {
		log.WithError(err).Error("failed to create new token")
	}

	return &IOSNotificationsService{
		APNSToken: token,
		IsActive:  true,
	}
}

func (s *CampusServer) GetIOSNotificationsService() *IOSNotificationsService {
	return s.iOSNotificationsService
}
