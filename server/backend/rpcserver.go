package backend

import (
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CampusServer struct {
	pb.UnimplementedCampusServer
	db        *gorm.DB
	deviceBuf *deviceBuffer // deviceBuf stores all devices from recent request and flushes them to db
}

// Verify that CampusServer implements the pb.CampusServer interface
var _ pb.CampusServer = (*CampusServer)(nil)

func New(db *gorm.DB) *CampusServer {
	log.Trace("Server starting up")
	return &CampusServer{
		db:        db,
		deviceBuf: newDeviceBuffer(),
	}
}
