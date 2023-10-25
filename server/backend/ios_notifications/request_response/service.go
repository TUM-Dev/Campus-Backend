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

var collectedNewGrades = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "ios_new_grades",
	Help:    "The total number of processed events",
	Buckets: prometheus.LinearBuckets(0, 5, 5),
})

func (service *Service) HandleDeviceRequestResponse(req *pb.IOSDeviceRequestResponseRequest, apnsIsActive bool) (*pb.IOSDeviceRequestResponseReply, error) {
	log.WithField("requestId", req.RequestId).Trace("Handling request")

	requestLog, err := service.Repository.GetIOSDeviceRequest(req.RequestId)
	if err != nil {
		log.WithError(err).Error("Could not get request")
		return nil, status.Error(codes.Internal, "Could not get request, probably request is already outdated")
	}

	switch requestLog.RequestType {
	case model.IOSBackgroundCampusTokenRequest.String():
		if req.Payload == "" {
			return nil, status.Error(codes.InvalidArgument, "Payload is empty")
		}

		if !apnsIsActive {
			return nil, status.Error(codes.Internal, "APNS is not active")
		}

		return service.handleDeviceCampusTokenRequest(requestLog, req.Payload)
	default:
		return nil, status.Error(codes.InvalidArgument, "Unknown request type")
	}
}

func (service *Service) handleDeviceCampusTokenRequest(requestLog *model.IOSDeviceRequestLog, campusToken string) (*pb.IOSDeviceRequestResponseReply, error) {
	log.WithField("DeviceID", requestLog.DeviceID).Trace("Handling campus token request")

	userRepo := device.NewRepository(service.Repository.DB)

	device, err := userRepo.GetDevice(requestLog.DeviceID)
	if err != nil {
		log.WithError(err).Error("Could not get device")
		return nil, status.Error(codes.Internal, "Could not get device")
	}

	apiGrades, err := campus_api.FetchGrades(campusToken)
	if err != nil {
		log.WithError(err).Error("Could not fetch grades")
		return nil, status.Error(codes.Internal, "Could not handle grades request")
	}

	oldEncryptedGrades, err := service.Repository.GetIOSEncryptedGrades(requestLog.DeviceID)
	if err != nil {
		log.WithError(err).Error("Could not get old grades")
		return nil, status.Error(codes.Internal, "Could not handle grades request")
	}

	oldGrades, err := decryptGrades(oldEncryptedGrades, campusToken)
	if err != nil {
		log.WithError(err).Error("Could not decrypt old grades")
		return nil, status.Error(codes.Internal, "Could not handle grades request")
	}

	newGrades := compareAndFindNewGrades(apiGrades.Grades, oldGrades)
	collectedNewGrades.Observe(float64(len(newGrades)))
	if len(newGrades) == 0 {
		log.Info("No new grades found")
		if err := service.Repository.DeleteAllRequestLogsForThisDeviceWithType(requestLog); err != nil {
			log.WithError(err).Error("Could not delete request logs")
		}
		return &pb.IOSDeviceRequestResponseReply{
			Message: "Successfully handled request",
		}, nil
	}

	err = service.Repository.DeleteEncryptedGrades(requestLog.DeviceID)
	if err != nil {
		log.WithError(err).Error("Could not delete old grades")
		return nil, status.Error(codes.Internal, "Could not handle grades request")
	}

	service.encryptGradesAndStoreInDatabase(apiGrades.Grades, requestLog.DeviceID, campusToken)

	log.WithFields(log.Fields{"old": len(oldGrades), "new": len(newGrades)}).Info("Found grades")

	if len(newGrades) > 0 && len(oldGrades) > 0 {
		apnsRepository := apns.NewRepository(service.Repository.DB, service.Repository.Token)
		sendGradesToDevice(device, newGrades, apnsRepository)
	}

	if err := service.Repository.DeleteAllRequestLogsForThisDeviceWithType(requestLog); err != nil {
		log.WithError(err).Error("Could not delete request logs")
	}

	return &pb.IOSDeviceRequestResponseReply{
		Message: "Successfully handled request",
	}, nil
}

func decryptGrades(grades []model.IOSEncryptedGrade, campusToken string) ([]model.IOSEncryptedGrade, error) {
	oldGrades := make([]model.IOSEncryptedGrade, len(grades))
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

func compareAndFindNewGrades(newGrades []model.IOSGrade, oldGrades []model.IOSEncryptedGrade) []model.IOSGrade {
	var grades []model.IOSGrade
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

func (service *Service) encryptGradesAndStoreInDatabase(grades []model.IOSGrade, deviceID string, campusToken string) {
	for _, grade := range grades {
		encryptedGrade := model.IOSEncryptedGrade{
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

func sendGradesToDevice(device *model.IOSDevice, grades []model.IOSGrade, apns *apns.Repository) {
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

	if _, err := apns.SendNotification(notificationPayload, model.IOSAPNSPushTypeAlert); err != nil {
		log.WithField("DeviceID", device.DeviceID).WithError(err).Error("Could not send notification")
	}
}

func NewService(repo *Repository) *Service {
	return &Service{
		Repository: repo,
	}
}
