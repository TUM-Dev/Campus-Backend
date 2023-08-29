package backend

import (
	"context"
	"database/sql"
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *CampusServer) SendFeedback(ctx context.Context, req *pb.SendFeedbackRequest) (*emptypb.Empty, error) {
	if err := s.checkDevice(ctx); err != nil {
		return nil, err
	}

	reciever := "app@tum.de"
	if req.Topic != "tca" {
		reciever = "kontakt@tum.de"
	}

	feedback := model.Feedback{
		EmailId:    sql.NullString{String: req.EmailId, Valid: true},
		Receiver:   sql.NullString{String: reciever, Valid: true},
		ReplyTo:    sql.NullString{String: req.Email, Valid: true},
		Feedback:   sql.NullString{String: req.Message, Valid: true},
		ImageCount: req.ImageCount,
		Latitude:   sql.NullFloat64{Float64: req.Latitude, Valid: true},
		Longitude:  sql.NullFloat64{Float64: req.Longitude, Valid: true},
		AppVersion: sql.NullString{String: req.AppVersion, Valid: true},
		OsVersion:  sql.NullString{String: req.OsVersion, Valid: true},
	}
	if err := s.db.Model(&model.Feedback{}).Create(&feedback).Error; err != nil {
		log.WithError(err).Error("Error while creating feedback")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
