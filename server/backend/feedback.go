package backend

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"slices"
	"strings"
	"time"

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

// CreateFeedback accepts a stream of feedback messages from the client and stores them in the database/file system.
func (s *CampusServer) CreateFeedback(stream pb.Campus_CreateFeedbackServer) error {
	// receive metadata
	id, err := uuid.NewGen().NewV4()
	if err != nil {
		log.WithError(err).Error("Error generating uuid")
		return status.Error(codes.Internal, "Error starting feedback submission")
	}
	feedback := &model.Feedback{EmailId: id.String(), Recipient: "app@tum.de"}

	// download images
	dbPath := path.Join("feedback", feedback.EmailId)
	var uploadedFilenames []*string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.WithError(err).Error("Error receiving feedback")
			deleteUploaded(feedback.EmailId)
			return status.Error(codes.Internal, "Error receiving feedback")
		}
		mergeFeedback(feedback, req)

		if len(req.Attachment) > 0 {
			filename := handleImageUpload(req.Attachment, len(uploadedFilenames), dbPath)
			if filename != nil {
				uploadedFilenames = append(uploadedFilenames, filename)
			}
		}
	}
	feedback.ImageCount = int32(len(uploadedFilenames))
	// validate feedback
	feedback.Feedback = strings.TrimSpace(feedback.Feedback)
	feedback.Feedback = strings.ReplaceAll(feedback.Feedback, "  ", " ")
	feedback.Feedback = strings.ToValidUTF8(feedback.Feedback, "?")
	if feedback.Feedback == "" && feedback.ImageCount == 0 {
		return status.Error(codes.InvalidArgument, "Please attach an image or feedback for us")
	}
	if feedback.ReplyToEmail.Valid {
		now := time.Now()
		fiveMinutesAgo := now.Add(time.Minute * -5).Unix()
		lastFeedback, feedbackExisted := s.feedbackEmailLastReuestAt.LoadOrStore(feedback.ReplyToEmail.String, now.Unix())
		if feedbackExisted && lastFeedback.(int64) >= fiveMinutesAgo {
			return status.Error(codes.ResourceExhausted, fmt.Sprintf("You have already send a feedback recently. Please wait %d seconds", lastFeedback.(int64)-fiveMinutesAgo))
		}
	}
	// save feedback to db
	if err := s.db.WithContext(stream.Context()).Transaction(func(tx *gorm.DB) error {
		var existingFeeedbackCnt int64
		if err := tx.Model(&feedback).Where("receiver=? AND reply_to_email=? AND feedback=? AND app_version=?", feedback.Recipient, feedback.ReplyToEmail, feedback.Feedback, feedback.AppVersion).Count(&existingFeeedbackCnt).Error; err != nil {
			return err
		}
		if existingFeeedbackCnt != 0 {
			return gorm.ErrDuplicatedKey
		}
		for _, filename := range uploadedFilenames {
			if err := tx.Create(&model.File{
				Name:       *filename,
				Path:       dbPath,
				Downloads:  1,
				Downloaded: null.BoolFrom(true),
			}).Error; err != nil {
				return err
			}
		}
		return tx.Create(feedback).Error
	}); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return status.Error(codes.AlreadyExists, "Feedback already exists")
		}
		log.WithError(err).Error("Error creating feedback")
		deleteUploaded(feedback.EmailId)
		return status.Error(codes.Internal, "Error creating feedback")
	}

	if err := stream.SendAndClose(&pb.CreateFeedbackReply{}); err != nil {
		log.WithError(err).Error("Error sending feedback-reply")
		return status.Error(codes.Internal, "Error sending feedback-reply")
	}
	return nil
}

// deleteUploaded deletes all uploaded images from the filesystem
func deleteUploaded(dbPath string) {
	if err := os.RemoveAll(cron.StorageDir + dbPath); err != nil {
		log.WithError(err).WithField("path", dbPath).Error("Error deleting uploaded images from filesystem")
	}
}

func handleImageUpload(content []byte, imageCounter int, dbPath string) *string {
	filename, realFilePath := inferFileName(mimetype.Detect(content), dbPath, imageCounter)
	if filename == nil {
		return nil // the filetype is not accepted by us
	}

	if err := os.MkdirAll(path.Dir(*realFilePath), 0755); err != nil {
		log.WithError(err).WithField("dbPath", dbPath).Error("Error creating directory for feedback")
		return nil
	}
	out, err := os.Create(*realFilePath)
	if err != nil {
		log.WithError(err).WithField("path", dbPath).Error("Error creating file for feedback")
		return nil
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.WithError(err).WithField("path", dbPath).Error("Error while closing file")
		}
	}(out)
	if _, err := io.Copy(out, bytes.NewReader(content)); err != nil {
		log.WithError(err).WithField("path", dbPath).Error("Error while writing file")
		if err := os.Remove(*realFilePath); err != nil {
			log.WithError(err).WithField("path", dbPath).Warn("Could not clean up file")
		}
		return nil
	}
	return filename
}

func inferFileName(mime *mimetype.MIME, dbPath string, counter int) (*string, *string) {
	allowedExt := []string{".jpeg", ".jpg", ".png", ".webp", ".md", ".txt", ".pdf"}
	if !slices.Contains(allowedExt, mime.Extension()) {
		return nil, nil
	}

	filename := fmt.Sprintf("%d%s", counter, mime.Extension())
	realFilePath := path.Join(cron.StorageDir, dbPath, filename)
	return &filename, &realFilePath
}

func mergeFeedback(feedback *model.Feedback, req *pb.CreateFeedbackRequest) {
	if req.Recipient.Enum() != nil {
		feedback.Recipient = receiverFromTopic(req.Recipient)
	}
	if req.OsVersion != "" {
		feedback.OsVersion = null.StringFrom(req.OsVersion)
	}
	if req.AppVersion != "" {
		feedback.AppVersion = null.StringFrom(req.AppVersion)
	}
	if req.Location != nil && req.Location.Longitude != 0 && req.Location.Latitude != 0 {
		feedback.Longitude = null.FloatFrom(req.Location.Longitude)
		feedback.Latitude = null.FloatFrom(req.Location.Latitude)
	}
	if req.Message != "" {
		feedback.Feedback = req.Message
	}
	if req.FromEmail != "" {
		feedback.ReplyToEmail = null.StringFrom(req.FromEmail)
	}
	if req.FromName != "" {
		feedback.ReplyToName = null.StringFrom(req.FromEmail)
	}
}

func receiverFromTopic(topic pb.CreateFeedbackRequest_Recipient) string {
	switch topic {
	case pb.CreateFeedbackRequest_TUM_DEV:
		return "app@tum.de"
	default:
		return "kontakt@tum.de"
	}
}
