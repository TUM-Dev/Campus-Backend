package backend

import (
	"context"
	"database/sql"
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"os"
	"path/filepath"
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

func (s *CampusServer) SendFeedbackImage(stream pb.Campus_SendFeedbackImageServer) error {
	if err := s.checkDevice(stream.Context()); err != nil {
		return err
	}

	// prepare temporary file to save steam to
	file, err := os.CreateTemp("", "*")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.WithError(err).Error("could not close close file for a feedback-image")
		}
	}(file)
	if err != nil {
		log.WithError(err).Error("failed to create temporary file for a feedback-image")
		return status.Error(codes.Internal, err.Error())
	}
	// download the file to a temporary file
	fileSize := 0
	finalDestination := ""
	for {
		req, err := stream.Recv()
		if finalDestination == "" {
			finalDestination = fmt.Sprintf("/Storage/feedback/%d/%d.png", req.GetId(), req.GetImageNr())
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fields := log.Fields{"finalDestination": finalDestination, "fileSize": fileSize}
			log.WithError(err).WithFields(fields).Error("failed to recive feedbackImage")
			return status.Error(codes.Aborted, err.Error())
		}
		bytesWritten, err := file.Write(req.GetContent())
		fileSize += bytesWritten
		log.Debug("received a chunk with size: ", bytesWritten)
		if err != nil {
			fields := log.Fields{"finalDestination": finalDestination, "fileSize": fileSize}
			log.WithError(err).WithFields(fields).Error("failed to write chunk for feedbackImage")
			return status.Error(codes.ResourceExhausted, err.Error())
		}
	}
	fields := log.Fields{"finalDestination": finalDestination, "fileSize": fileSize, "temporaryFilename": file.Name()}
	log.WithFields(fields).Debug("received feedbackImage without issue")

	// move the file to a different directory
	parent := filepath.Dir(finalDestination)
	if err = os.MkdirAll(parent, os.ModePerm); err != nil {
		log.WithError(err).WithField("finalDestination", finalDestination).Error("failed to make all directories")
		return err
	}
	if err := os.Rename(file.Name(), finalDestination); err != nil {
		fields := log.Fields{"finalDestination": finalDestination, "temporaryFilename": file.Name()}
		log.WithError(err).WithFields(fields).Error("could not move file to the required location")
		return status.Error(codes.Internal, err.Error())
	}
	log.WithFields(fields).Debug("moved the temporary feedbackImage to the correct location")

	return nil
}
