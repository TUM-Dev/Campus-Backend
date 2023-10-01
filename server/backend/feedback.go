package backend

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gofrs/uuid/v5"
	"github.com/guregu/null"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// NewFeedback accepts a stream of feedback messages from the client and stores them in the database/file system.
func (s *CampusServer) NewFeedback(stream pb.Campus_NewFeedbackServer) error {
	return s.db.WithContext(stream.Context()).Transaction(func(tx *gorm.DB) error {
		// receive metadata
		imageCount := int32(0)
		id, err := uuid.NewGen().NewV7()
		if err != nil {
			log.WithError(err).Error("Error generating uuid")
			return status.Error(codes.Internal, "Error starting feedback submission")
		}
		feedback := &model.Feedback{EmailId: null.StringFrom(id.String())}

		// download images
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.WithError(err).Error("Error receiving feedback")
				deleteUploaded(tx, feedback.EmailId.String)
				return status.Error(codes.Internal, "Error receiving feedback")
			}
			mergeFeedback(feedback, req)

			if len(req.Attachment) > 0 {
				imageCount = handleImageUpload(tx, &req.Attachment, imageCount, feedback.EmailId.String)
			}
		}
		feedback.ImageCount = imageCount
		if err := tx.Create(feedback).Error; err != nil {
			log.WithError(err).Error("Error creating feedback")
			return status.Error(codes.Internal, "Error creating feedback")
		}
		if err := stream.SendAndClose(&pb.NewFeedbackReply{}); err != nil {
			log.WithError(err).Error("Error sending feedbackreply")
			return status.Error(codes.Internal, "Error sending feedbackreply")
		}
		return nil
	})
}

func deleteUploaded(tx *gorm.DB, id string) {
	// delete uploaded images
	dbPath := fmt.Sprintf("feedback/%s", id)
	if err := tx.Find(&model.File{Path: dbPath}).Delete(&model.File{}).Error; err != nil {
		log.WithError(err).WithField("path", dbPath).Error("Error deleting uploaded images from db")
	}
	if err := os.RemoveAll(cron.StorageDir + dbPath); err != nil {
		log.WithError(err).WithField("path", dbPath).Error("Error deleting uploaded images from filesystem")
	}
}

func handleImageUpload(tx *gorm.DB, content *[]byte, imageCounter int32, id string) int32 {
	filename := inferFileName(content, imageCounter)
	dbPath := path.Join("feedback", id)
	realFilePath := path.Join(cron.StorageDir, dbPath, filename)

	if err := os.MkdirAll(path.Dir(realFilePath), 0755); err != nil {
		log.WithError(err).WithField("dbPath", dbPath).Error("Error creating directory for feedback")
		return imageCounter
	}
	out, err := os.Create(realFilePath)
	if err != nil {
		log.WithError(err).WithField("path", dbPath).Error("Error creating file for feedback")
		return imageCounter
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.WithError(err).WithField("path", dbPath).Error("Error while closing file")
		}
	}(out)
	if _, err := io.Copy(out, bytes.NewReader(*content)); err != nil {
		log.WithError(err).WithField("path", dbPath).Error("Error while writing file")
		if err := os.Remove(realFilePath); err != nil {
			log.WithError(err).WithField("path", dbPath).Warn("Could not clean up file")
		}
		return imageCounter
	}

	tx.Create(&model.File{
		Name:       filename,
		Path:       dbPath,
		Downloads:  1,
		Downloaded: null.BoolFrom(true),
	})
	return imageCounter + 1
}

func inferFileName(content *[]byte, counter int32) string {
	ext := mimetype.Detect(*content).Extension()
	return fmt.Sprintf("%d%s", counter, ext)
}

func mergeFeedback(feedback *model.Feedback, req *pb.NewFeedbackRequest) {
	if req.Recipient.Enum() != nil {
		feedback.Recipient = null.StringFrom(receiverFromTopic(req.Recipient))
	}
	if req.Metadata != nil {
		feedback.OsVersion = null.StringFrom(req.Metadata.OsVersion)
		feedback.AppVersion = null.StringFrom(req.Metadata.AppVersion)
		feedback.Longitude = null.FloatFrom(req.Metadata.Longitude)
		feedback.Latitude = null.FloatFrom(req.Metadata.Latitude)
	}
	if req.Message != "" {
		feedback.Feedback = null.StringFrom(req.Message)
	}
	if req.FromEmail != "" {
		feedback.ReplyTo = null.StringFrom(req.FromEmail)
	}
}

func receiverFromTopic(topic pb.NewFeedbackRequest_Recipient) string {
	switch topic {
	case pb.NewFeedbackRequest_TUM_DEV:
		return "app@tum.de"
	default:
		return "kontakt@tum.de"
	}
}
