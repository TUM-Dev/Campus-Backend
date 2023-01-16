package ios_request_response

import (
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/api"
	"github.com/TUM-Dev/Campus-Backend/backend/campus_api"
	"github.com/TUM-Dev/Campus-Backend/backend/ios_notifications/ios_apns"
	"github.com/TUM-Dev/Campus-Backend/backend/ios_notifications/ios_logging"
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Repository *Repository
	Logger     *ios_logging.Service
}

func (service *Service) HandleDeviceRequestResponse(request *pb.IOSDeviceRequestResponseRequest) (*pb.IOSDeviceRequestResponseReply, error) {

	log.Infof("Received request with id %s with payload: %s", request.GetRequestId(), request.GetPayload())

	requestId := request.GetRequestId()

	service.Logger.LogTokenRequest("Received Token Request: %s", requestId)

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

	oldGrades, err := service.decryptGrades(oldEncryptedGrades, campusToken)

	if err != nil {
		return nil, err
	}

	// compare old and new grades
	newGrades := service.compareAndFindNewGrades(apiGrades.Grades, oldGrades)

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

	service.encryptGradesAndStoreInDatabase(apiGrades.Grades, requestLog.DeviceID, campusToken)

	if len(newGrades) > 0 {
		apnsRepository := ios_apns.NewRepository(service.Repository.DB, service.Repository.Token)
		service.sendGradesToDevice(requestLog.DeviceID, newGrades, apnsRepository)
	}

	err = service.Repository.DeleteRequestLog(requestLog.RequestID)

	if err != nil {
		log.Error("Could not delete request log: ", err)
	}

	return &pb.IOSDeviceRequestResponseReply{
		Message: "Successfully handled request",
	}, nil
}

func (service *Service) decryptGrades(grades []model.IOSEncryptedGrade, campusToken string) ([]model.IOSEncryptedGrade, error) {
	oldGrades := make([]model.IOSEncryptedGrade, len(grades))
	for i, encryptedGrade := range grades {
		err := encryptedGrade.Decrypt(campusToken)

		if err != nil {
			log.Error("Could not decrypt grade: ", err)
			return nil, status.Error(codes.Internal, "Could not decrypt grade")
		}

		oldGrades[i] = encryptedGrade
	}

	return oldGrades, nil
}

func (service *Service) compareAndFindNewGrades(newGrades []model.IOSGrade, oldGrades []model.IOSEncryptedGrade) []model.IOSGrade {
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
			newGrades = append(newGrades, grade)
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
			log.Error("Could not encrypt grade: ", err)
		}

		err = service.Repository.SaveEncryptedGrade(&encryptedGrade)

		if err != nil {
			log.Error("Could not save grade: ", err)
		}
	}
}

func (service *Service) sendGradesToDevice(deviceId string, grades []model.IOSGrade, apns *ios_apns.Repository) {
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

	notificationPayload := model.NewIOSNotificationPayload(deviceId).Alert(alertTitle, "", alertBody)

	_, err := apns.SendAlertNotification(notificationPayload)

	if err != nil {
		log.Error("Could not send notification: ", err)
	}
}

func NewService(repo *Repository) *Service {
	return &Service{
		Repository: repo,
		Logger:     ios_logging.NewLogger(repo.DB),
	}
}
