// Package request_response provides functionality to handle device requests.
// Device Requests are requests that are sent from the device to the server when the
// device received a background push notification from the backend.
package request_response

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/backend/campus_api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/apns"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/device"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Repository *Repository
}

var (
	ErrOutdatedRequest      = status.Error(codes.Internal, "Could not get request, probably request is already outdated")
	ErrEmptyPayload         = status.Error(codes.InvalidArgument, "Payload is empty")
	ErrUnknownRequestType   = status.Error(codes.InvalidArgument, "Unknown request type")
	ErrInternalHandleGrades = status.Error(codes.Internal, "Could not handle grades request")
	ErrCouldNotGetDevice    = status.Error(codes.Internal, "Could not get device")
	ErrAPNSNotActive        = status.Error(codes.Internal, "APNS is not active")

	collectedNewGrades = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "ios_new_grades",
		Help:    "The total number of processed events",
		Buckets: prometheus.LinearBuckets(0, 5, 5),
	})
)

func (service *Service) HandleDeviceRequestResponse(request *pb.IOSDeviceRequestResponseRequest, apnsIsActive bool) (*pb.IOSDeviceRequestResponseReply, error) {
	// requestId refers to the request id that was sent to the device and stored in the Database
	requestId := request.GetRequestId()

	log.WithField("requestId", requestId).Info("Handling request")

	requestLog, err := service.Repository.GetIOSDeviceRequest(requestId)

	if err != nil {
		log.WithError(err).Error("Could not get request")
		return nil, ErrOutdatedRequest
	}

	switch requestLog.RequestType {
	case model.IOSBackgroundCampusTokenRequest.String():
		campusToken := request.GetPayload()

		if campusToken == "" {
			return nil, ErrEmptyPayload
		}

		if !apnsIsActive {
			return nil, ErrAPNSNotActive
		}

		return service.handleDeviceCampusTokenRequest(requestLog, campusToken)
	default:
		return nil, ErrUnknownRequestType
	}
}

func (service *Service) handleDeviceCampusTokenRequest(requestLog *model.IOSDeviceRequestLog, campusToken string) (*pb.IOSDeviceRequestResponseReply, error) {
	log.WithField("DeviceID", requestLog.DeviceID).Info("Handling campus token request")

	userRepo := device.NewRepository(service.Repository.DB)

	device, err := userRepo.GetDevice(requestLog.DeviceID)

	if err != nil {
		log.WithError(err).Error("Could not get device")
		return nil, ErrCouldNotGetDevice
	}

	apiGrades, err := campus_api.FetchGrades(campusToken)
	if err != nil {
		log.WithError(err).Error("Could not fetch grades")
		return nil, ErrInternalHandleGrades
	}

	oldEncryptedGrades, err := service.Repository.GetIOSEncryptedGrades(requestLog.DeviceID)
	if err != nil {
		log.WithError(err).Error("Could not get old grades")
		return nil, ErrInternalHandleGrades
	}

	oldGrades, err := decryptGrades(oldEncryptedGrades, campusToken)
	if err != nil {
		log.WithError(err).Error("Could not decrypt old grades")
		return nil, ErrInternalHandleGrades
	}

	newGrades := compareAndFindNewGrades(apiGrades.Grades, oldGrades)
	collectedNewGrades.Observe(float64(len(newGrades)))
	if len(newGrades) == 0 {
		log.Info("No new grades found")
		service.deleteRequestLog(requestLog)

		if err != nil {
			log.Error("Could not send push notification: ", err)
			return nil, ErrInternalHandleGrades
		}

		return &pb.IOSDeviceRequestResponseReply{
			Message: "Successfully handled request",
		}, nil
	}

	err = service.Repository.DeleteEncryptedGrades(requestLog.DeviceID)

	if err != nil {
		log.WithError(err).Error("Could not delete old grades")
		return nil, ErrInternalHandleGrades
	}

	service.encryptGradesAndStoreInDatabase(apiGrades.Grades, requestLog.DeviceID, campusToken)

	log.WithFields(log.Fields{"old": len(oldGrades), "new": len(newGrades)}).Info("Found grades")

	if len(newGrades) > 0 && len(oldGrades) > 0 {
		apnsRepository := apns.NewRepository(service.Repository.DB, service.Repository.Token)
		sendGradesToDevice(device, newGrades, apnsRepository)
	}

	service.deleteRequestLog(requestLog)

	return &pb.IOSDeviceRequestResponseReply{
		Message: "Successfully handled request",
	}, nil
}

func (service *Service) deleteRequestLog(requestLog *model.IOSDeviceRequestLog) {
	err := service.Repository.DeleteAllRequestLogsForThisDeviceWithType(requestLog)

	if err != nil {
		log.WithError(err).Error("Could not delete request logs")
	}
}

func decryptGrades(grades []model.EncryptedGrade, campusToken string) ([]model.EncryptedGrade, error) {
	oldGrades := make([]model.EncryptedGrade, len(grades))
	for i, encryptedGrade := range grades {
		err := encryptedGrade.Decrypt(campusToken)

		if err != nil {
			log.WithError(err).Error("Could not decrypt grade")
			return nil, status.Error(codes.Internal, "Could not decrypt grade")
		}

		oldGrades[i] = encryptedGrade
	}

	return oldGrades, nil
}

func compareAndFindNewGrades(newGrades []model.Grade, oldGrades []model.EncryptedGrade) []model.Grade {
	var grades []model.Grade
	for _, grade := range newGrades {
		found := false
		for _, oldGrade := range oldGrades {
			if grade.CompareToEncrypted(&oldGrade) {
				found = true
				break
			}
		}

		if !found {
			grades = append(grades, grade)
		}
	}

	return grades
}

func (service *Service) encryptGradesAndStoreInDatabase(grades []model.Grade, deviceID string, campusToken string) {
	for _, grade := range grades {
		encryptedGrade := model.EncryptedGrade{
			Grade:        grade.Grade,
			DeviceID:     deviceID,
			LectureTitle: grade.LectureTitle,
		}

		err := encryptedGrade.Encrypt(campusToken)
		if err != nil {
			log.WithError(err).Error("Could not encrypt grade")
		}

		err = service.Repository.SaveEncryptedGrade(&encryptedGrade)
		if err != nil {
			log.WithError(err).Error("Could not save grade")
		}
	}
}

func sendGradesToDevice(device *model.IOSDevice, grades []model.Grade, apns *apns.Repository) {
	alertTitle := fmt.Sprintf("%d New Grades Available", len(grades))

	if len(grades) == 1 {
		alertTitle = "New Grade Available"
	}

	var alertBody string
	for i, grade := range grades {
		if i == 0 {
			alertBody = fmt.Sprintf("%s: %s", grade.LectureTitle, grade.Grade)
		} else {
			alertBody = fmt.Sprintf("%s\n %s: %s", alertBody, grade.LectureTitle, grade.Grade)
		}
	}

	notificationPayload := model.NewIOSNotificationPayload(device.DeviceID).
		Alert(alertTitle, "", alertBody).
		Encrypt(device.PublicKey)

	log.WithField("DeviceID", device.DeviceID).Info("Sending push notification")

	_, err := apns.SendAlertNotification(notificationPayload)
	if err != nil {
		log.WithField("DeviceID", device.DeviceID).WithError(err).Error("Could not send notification")
	}
}

func NewService(repo *Repository) *Service {
	return &Service{
		Repository: repo,
	}
}
