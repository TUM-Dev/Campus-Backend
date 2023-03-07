// Package ios_request_response provides functionality to handle device requests.
// Device Requests are requests that are sent from the device to the server when the
// device received a background push notification from the backend.
package ios_request_response

import (
	"fmt"
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/campus_api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/influx"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_grades"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_lectures"
	"github.com/TUM-Dev/Campus-Backend/server/backend/utils"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Repository *Repository
	APNs       *ios_apns.Service
}

var (
	ErrOutdatedRequest       = status.Error(codes.Internal, "Could not get request, probably request is already outdated")
	ErrEmptyPayload          = status.Error(codes.InvalidArgument, "Payload is empty")
	ErrUnknownRequestType    = status.Error(codes.InvalidArgument, "Unknown request type")
	ErrInternalHandleGrades  = status.Error(codes.Internal, "Could not handle grades request")
	ErrCouldNotGetDevice     = status.Error(codes.Internal, "Could not get device")
	ErrAPNSNotActive         = status.Error(codes.Internal, "APNS is not active")
	ErrCouldNotDecryptGrades = status.Error(codes.Internal, "Could not decrypt grades")
)

func (service *Service) HandleDeviceRequestResponse(request *pb.IOSDeviceRequestResponseRequest, apnsIsActive bool) (*pb.IOSDeviceRequestResponseReply, error) {
	// requestId refers to the request id that was sent to the device and stored in the Database
	requestId := request.GetRequestId()

	log.Infof("Handling request with id %s", requestId)

	requestLog, err := service.Repository.GetIOSDeviceRequest(requestId)
	if err != nil {
		log.WithError(err).Error("Could not get request")
		return nil, ErrOutdatedRequest
	}

	if requestLog.HandledAt.Valid {
		log.Warnf("Request with id %s is already handled", requestId)
		return nil, ErrOutdatedRequest
	}

	influx.LogIOSBackgroundRequestResponse(requestLog.DeviceID, requestLog.RequestType)

	service.setRequestLogAsHandled(requestLog)

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
	case model.IOSBackgroundLectureUpdateRequest.String():
		campusToken := request.GetPayload()

		if campusToken == "" {
			return nil, ErrEmptyPayload
		}

		if !apnsIsActive {
			return nil, ErrAPNSNotActive
		}

		return service.handleUpdateLecturesRequest(requestLog, campusToken)
	default:
		return nil, ErrUnknownRequestType
	}
}

func (service *Service) handleUpdateLecturesRequest(requestLog *model.IOSDeviceRequestLog, campusToken string) (*pb.IOSDeviceRequestResponseReply, error) {
	log.Infof("Handling update lectures request for device %s", requestLog.DeviceID)

	devicesRepo := ios_device.NewRepository(service.Repository.DB)
	device, err := devicesRepo.GetDevice(requestLog.DeviceID)
	if err != nil {
		log.WithError(err).Error("Could not get device")
		return nil, ErrCouldNotGetDevice
	}

	newGrades, err := service.findGetAndUpdateNewGrades(requestLog.DeviceID, campusToken, requestLog.RequestID, true)
	if err != nil {
		log.WithError(err).Error("Could not get new grades")
		return nil, err
	}

	if len(*newGrades) > 0 {
		apnsRepository := ios_apns.NewRepository(service.Repository.DB, service.Repository.Token)
		sendGradesToDevice(device, newGrades, apnsRepository)

		influx.LogIOSNewGrades(requestLog.DeviceID, len(*newGrades))

		service.notifyOtherDevicesAboutNewGrades(device.DeviceID, newGrades)
	}

	return &pb.IOSDeviceRequestResponseReply{
		Message: "Successfully handled request",
	}, nil
}

func (service *Service) handleDeviceCampusTokenRequest(requestLog *model.IOSDeviceRequestLog, campusToken string) (*pb.IOSDeviceRequestResponseReply, error) {
	log.Infof("Handling campus token request for device %s", requestLog.DeviceID)

	devicesRepo := ios_device.NewRepository(service.Repository.DB)
	device, err := devicesRepo.GetDevice(requestLog.DeviceID)
	if err != nil {
		log.WithError(err).Error("Could not get device")
		return nil, ErrCouldNotGetDevice
	}

	newGrades, err := service.findGetAndUpdateNewGrades(requestLog.DeviceID, campusToken, requestLog.RequestID, false)
	if err != nil {
		log.WithError(err).Error("Could not get new grades")
		return nil, err
	}

	if len(*newGrades) > 0 {
		apnsRepository := ios_apns.NewRepository(service.Repository.DB, service.Repository.Token)
		sendGradesToDevice(device, newGrades, apnsRepository)

		influx.LogIOSNewGrades(requestLog.DeviceID, len(*newGrades))
	}

	return &pb.IOSDeviceRequestResponseReply{
		Message: "Successfully handled request",
	}, nil
}

func (service *Service) notifyOtherDevicesAboutNewGrades(deviceId string, newGrades *[]model.Grade) {
	log.Infof("Notifying other devices about new grades")
	lecturesRepo := ios_lectures.NewRepository(service.Repository.DB)

	for _, grade := range *newGrades {
		devicesThatAttendLecture, err := lecturesRepo.FindOtherDevicesThatAttendLecture(grade.LectureTitle)
		log.Infof("Found %d devices that attend lecture %s", len(*devicesThatAttendLecture), grade.LectureTitle)
		if err != nil {
			log.WithError(err).Error("Could not find other devices that attend lecture")
			continue
		}

		utils.RunTasksInRoutines(devicesThatAttendLecture, func(device model.IOSDevice) {
			/* if device.DeviceID == deviceId {
				return
			}*/
			log.Infof("Notifying device %s about new grades", device.DeviceID)
			service.APNs.RequestGradeUpdateForDevice(device.DeviceID)
		}, 10)
	}
}

func (service *Service) findGetAndUpdateNewGrades(deviceId, campusToken, requestId string, updateLectures bool) (*[]model.Grade, error) {
	lecturesRepo := ios_lectures.NewRepository(service.Repository.DB)
	gradesRepo := ios_grades.NewRepository(service.Repository.DB)

	apiGrades, err := campus_api.FetchGrades(campusToken)
	if err != nil {
		log.Error("Could not fetch grades: ", err)
		return nil, ErrInternalHandleGrades
	}

	if updateLectures {
		lecturesRepo.SetLecturesLastUpdatedBy(requestId, &apiGrades.Grades)
	}

	oldGrades, err := gradesRepo.GetAndDecryptGrades(deviceId, campusToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "Could not decrypt grade")
	}

	newGrades := compareAndFindNewGrades(apiGrades.Grades, oldGrades)
	if len(newGrades) == 0 {
		log.Info("No new grades found")
		return &newGrades, nil
	}

	err = gradesRepo.DeleteEncryptedGrades(deviceId)
	if err != nil {
		log.Error("Could not delete old grades: ", err)
		return nil, ErrInternalHandleGrades
	}

	gradesRepo.EncryptAndSaveGrades(apiGrades.Grades, deviceId, campusToken)

	log.Infof("Found %d old grades and %d new grades", len(oldGrades), len(newGrades))

	// check if this is the first time the user has grades => no need to send notification
	if len(newGrades) > 0 && len(oldGrades) == 0 {
		return &([]model.Grade{}), nil
	}

	return &newGrades, nil
}

func (service *Service) setRequestLogAsHandled(requestLog *model.IOSDeviceRequestLog) {
	err := service.Repository.SetRequestLogAsHandled(requestLog)

	if err != nil {
		log.Error("Could not delete request logs: ", err)
	}
}

func compareAndFindNewGrades(newGrades []model.Grade, oldGrades []model.IOSEncryptedGrade) []model.Grade {
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

func sendGradesToDevice(device *model.IOSDevice, grades *[]model.Grade, apns *ios_apns.Repository) {
	alertTitle := fmt.Sprintf("%d New Grades Available", len(*grades))

	if len(*grades) == 1 {
		alertTitle = "New Grade Available"
	}

	var alertBody string
	for i, grade := range *grades {
		if i == 0 {
			alertBody = fmt.Sprintf("%s: %s", grade.LectureTitle, grade.Grade)
		} else {
			alertBody = fmt.Sprintf("%s\n %s: %s", alertBody, grade.LectureTitle, grade.Grade)
		}
	}

	notificationPayload := model.NewIOSNotificationPayload(device.DeviceID).
		Alert(alertTitle, "", alertBody).
		Encrypt(device.PublicKey)

	log.Infof("Sending push notification to device %s", device.DeviceID)

	_, err := apns.SendAlertNotification(notificationPayload)

	if err != nil {
		log.Error("Could not send notification: ", err)
	}
}

func NewService(repo *Repository, apns *ios_apns.Service) *Service {
	return &Service{
		Repository: repo,
		APNs:       apns,
	}
}
