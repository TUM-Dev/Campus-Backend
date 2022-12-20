package ios_request_response

import (
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

	if requestLog.RequestType != model.IOSBackgroundCampusTokenRequest.String() {
		return nil, status.Error(codes.Internal, "Request type is not implemented yet")
	}

	campusToken := request.GetPayload()

	if campusToken == "" {
		return nil, status.Error(codes.InvalidArgument, "Payload is empty")
	}

	grades, err := campus_api.FetchGrades(campusToken)

	if err != nil {
		log.Error("Could not fetch grades: ", err)
		return nil, status.Error(codes.Internal, "Could not fetch grades")
	}

	log.Infof("Fetched grades: %s", grades)

	for _, grade := range grades.Grades {
		encryptedGrade := model.IOSEncryptedGrade{
			Grade:        grade.Grade,
			DeviceID:     requestLog.DeviceID,
			LectureTitle: grade.LectureTitle,
		}

		encryptedGrade.Encrypt(campusToken)

		if err != nil {
			log.Error("Could not encrypt grade: ", err)
			return nil, status.Error(codes.Internal, "Could not encrypt grade")
		}

		err = service.Repository.SaveEncryptedGrade(&encryptedGrade)

		if err != nil {
			log.Error("Could not save grade: ", err)
			return nil, status.Error(codes.Internal, "Could not save grade")
		}
	}

	encryptedGrades, err := service.Repository.GetIOSEncryptedGrades(requestLog.DeviceID)

	if err != nil {
		log.Error("Could not get encrypted grades: ", err)
		return nil, status.Error(codes.Internal, "Could not get encrypted grades")
	}

	apnsRepository := ios_apns.NewRepository(service.Repository.DB, service.Repository.Token)

	alertTitle := "No New Grades Available"
	alertSubtitle := "You have no new grades available"

	if len(encryptedGrades) > 0 {
		grade := encryptedGrades[0]

		err := grade.Decrypt(campusToken)

		if err != nil {
			log.Error("Could not decrypt grade: ", err)
			return nil, status.Error(codes.Internal, "Could not decrypt grade")
		}

		alertTitle = "New Grades Available"
		alertSubtitle = grade.LectureTitle
	}

	notificationPayload := model.NewIOSNotificationPayload(requestLog.DeviceID).Alert(alertTitle, "", alertSubtitle)

	apnsRepository.SendAlertNotification(notificationPayload)

	if err != nil {
		log.Error("Could not get encrypted grades: ", err)
		return nil, status.Error(codes.Internal, "Could not get encrypted grades")
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
