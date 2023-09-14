package backend

import (
	"context"
	"errors"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *CampusServer) GetUpdateNote(ctx context.Context, req *pb.GetUpdateNoteRequest) (*pb.GetUpdateNoteReply, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	res := model.UpdateNote{VersionCode: req.Version}
	if err := s.db.First(&res).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(404, "No update note found")
	} else if err != nil {
		log.WithField("VersionCode", req.Version).WithError(err).Error("Failed to get update note")
		return nil, status.Error(500, "Internal server error")
	}

	return &pb.GetUpdateNoteReply{Message: res.Message, VersionName: res.VersionName}, nil
}
