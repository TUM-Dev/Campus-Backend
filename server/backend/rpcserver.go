package backend

import (
	"net"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns/ios_apns_jwt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func (s *CampusServer) GRPCServe(l net.Listener) error {
	grpcServer := grpc.NewServer()
	pb.RegisterCampusServer(grpcServer, s)
	if err := grpcServer.Serve(l); err != nil {
		log.WithError(err).Fatal("failed to serve")
	}
	return grpcServer.Serve(l)
}

type CampusServer struct {
	pb.UnimplementedCampusServer
	db                      *gorm.DB
	deviceBuf               *deviceBuffer // deviceBuf stores all devices from recent request and flushes them to db
	iOSNotificationsService *IOSNotificationsService
}

// Verify that CampusServer implements the pb.CampusServer interface
var _ pb.CampusServer = (*CampusServer)(nil)

func New(db *gorm.DB) *CampusServer {
	log.Trace("Server starting up")
	initTagRatingOptions(db)

	return &CampusServer{
		db:                      db,
		deviceBuf:               newDeviceBuffer(),
		iOSNotificationsService: NewIOSNotificationsService(),
	}
}

func NewIOSNotificationsService() *IOSNotificationsService {
	if err := ios_apns.ValidateRequirementsForIOSNotificationsService(); err != nil {
		log.WithError(err).Warn("failed to validate requirements for ios notifications service")

		return &IOSNotificationsService{
			APNSToken: nil,
			IsActive:  false,
		}
	}

	token, err := ios_apns_jwt.NewToken()
	if err != nil {
		log.WithError(err).Fatal("failed to create new token")
	}

	return &IOSNotificationsService{
		APNSToken: token,
		IsActive:  true,
	}
}

func (s *CampusServer) GetIOSNotificationsService() *IOSNotificationsService {
	return s.iOSNotificationsService
}
