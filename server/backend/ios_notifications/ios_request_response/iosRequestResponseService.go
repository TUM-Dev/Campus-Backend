package ios_request_response

import (
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"github.com/TUM-Dev/Campus-Backend/backend/campus_api"
	"github.com/TUM-Dev/Campus-Backend/backend/ios_notifications/ios_apns"
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Repository *Repository
}

func (service *Service) HandleDeviceRequestResponse(request *pb.IOSDeviceRequestResponseRequest) (*pb.IOSDeviceRequestResponseReply, error) {

	log.Infof("Received request with id %s with payload: %s", request.GetRequestId(), request.GetPayload())

	requestId := request.GetRequestId()

	requestLog, err := service.Repository.GetIOSDeviceRequest(requestId)

	if err != nil {
		return nil, status.Error(codes.Internal, "Could not get request, probably request is already outdated")
	}

	switch requestLog.RequestType {
	case model.IOSBackgroundCampusTokenRequest.String():
		campusToken := request.GetPayload()

		if campusToken == "" {
			return nil, status.Error(codes.InvalidArgument, "Payload is empty")
		}

		return service.handleDeviceCampusTokenRequest(requestLog, request.GetPayload())
	default:
		return nil, status.Error(codes.InvalidArgument, "Request type is not yet implemented")
	}
}

func (service *Service) handleDeviceCampusTokenRequest(requestLog *model.IOSDeviceRequestLog, campusToken string) (*pb.IOSDeviceRequestResponseReply, error) {
	apiGrades, err := campus_api.FetchGrades(campusToken)

	if err != nil {
		log.Error("Could not fetch grades: ", err)
		return nil, status.Error(codes.Internal, "Could not fetch grades")
	}

	oldEncryptedGrades, err := service.Repository.GetIOSEncryptedGrades(requestLog.DeviceID)

	if err != nil {
		log.Error("Could not get old grades: ", err)
		return nil, status.Error(codes.Internal, "Could not get old grades")
	}

	// decrypt old grades from database to compare them with the new ones
	oldGrades := make([]model.IOSEncryptedGrade, len(oldEncryptedGrades))
	for i, encryptedGrade := range oldEncryptedGrades {
		err := encryptedGrade.Decrypt(campusToken)

		if err != nil {
			log.Error("Could not decrypt grade: ", err)
			return nil, status.Error(codes.Internal, "Could not decrypt grade")
		}

		oldGrades[i] = encryptedGrade
	}

	// compare old and new grades
	var newGrades []model.IOSGrade
	for _, grade := range apiGrades.Grades {
		found := false
		for _, oldGrade := range oldGrades {
			if grade.CompareToEncrypted(&oldGrade) {
				found = true
				break
			}
		}

		if !found {
			newGrades = append(newGrades, grade)
		}
	}

	if len(newGrades) == 0 {
		return &pb.IOSDeviceRequestResponseReply{
			Message: "Successfully handled request",
		}, nil
	}

	err = service.Repository.DeleteEncryptedGrades(requestLog.DeviceID)

	if err != nil {
		log.Error("Could not delete old grades: ", err)
		return nil, status.Error(codes.Internal, "Could not delete old grades")
	}

	// encrypt new grades and save them to database
	for _, grade := range apiGrades.Grades {
		encryptedGrade := model.IOSEncryptedGrade{
			Grade:        grade.Grade,
			DeviceID:     requestLog.DeviceID,
			LectureTitle: grade.LectureTitle,
		}

		encryptedGrade.Encrypt(campusToken)

		if err != nil {
			log.Error("Could not encrypt grade: ", err)
		}

		err = service.Repository.SaveEncryptedGrade(&encryptedGrade)

		if err != nil {
			log.Error("Could not save grade: ", err)
		}
	}

	// send push notification to device
	apnsRepository := ios_apns.NewRepository(service.Repository.DB, service.Repository.Token)

	if len(newGrades) > 0 {
		alertTitle := fmt.Sprintf("%d New Grades Available", len(newGrades))

		if len(newGrades) == 1 {
			alertTitle = "New Grade Available"
		}

		var alertBody string
		for i, grade := range newGrades {
			if i == 0 {
				alertBody = fmt.Sprintf("%s: %s", grade.LectureTitle, grade.Grade)
			} else {
				alertBody = fmt.Sprintf("%s\n %s: %s", alertBody, grade.LectureTitle, grade.Grade)
			}
		}

		notificationPayload := model.NewIOSNotificationPayload(requestLog.DeviceID).Alert(alertTitle, "", alertBody)

		notification, err := apnsRepository.SendAlertNotification(notificationPayload)

		log.Infof("Sent notification with reason %s", notification.Reason)

		if err != nil {
			log.Error("Could not send notification: ", err)
		}
	}

	return &pb.IOSDeviceRequestResponseReply{
		Message: "Successfully handled request",
	}, nil
}

func NewService(db *Repository) *Service {
	return &Service{
		Repository: db,
	}
}
