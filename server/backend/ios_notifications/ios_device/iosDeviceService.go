// Package ios_device provides functions to register and remove ios devices
package ios_device

import (
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/campus_api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/influx"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_exams"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_grades"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Repository *Repository
}

var (
	ErrCouldNotRegisterDevice = status.Error(codes.Internal, "Could not register device")
	ErrCouldNotRemoveDevice   = status.Error(codes.Internal, "Could not remove device")
)

func (service *Service) RegisterDevice(request *pb.RegisterDeviceRequest) (*pb.RegisterDeviceReply, error) {
	device := model.IOSDevice{
		DeviceID:  request.GetDeviceId(),
		PublicKey: request.GetPublicKey(),
	}

	deviceAlreadyExisted, err := service.Repository.RegisterDevice(&device)
	if err != nil {
		return nil, ErrCouldNotRegisterDevice
	}

	if !deviceAlreadyExisted {
		service.handleFirstDeviceRegistration(request)
	}

	influx.LogIOSRegisterDevice(request.GetDeviceId())

	return &pb.RegisterDeviceReply{
		DeviceId: device.DeviceID,
	}, nil
}

func (service *Service) RemoveDevice(request *pb.RemoveDeviceRequest) (*pb.RemoveDeviceReply, error) {
	err := service.Repository.RemoveDevice(request.GetDeviceId())

	if err != nil {
		return nil, ErrCouldNotRemoveDevice
	}

	influx.LogIOSRemoveDevice(request.GetDeviceId())

	return &pb.RemoveDeviceReply{
		DeviceId: request.GetDeviceId(),
	}, nil
}

func (service *Service) handleFirstDeviceRegistration(request *pb.RegisterDeviceRequest) {
	campusToken := request.GetCampusApiToken()
	deviceId := request.GetDeviceId()

	influx.LogIOSRegisterDevice(deviceId)

	err := service.fetchAndStoreDeviceExams(deviceId, campusToken)
	if err != nil {
		log.WithError(err).Error("Could not fetch lectures of device")
	}

	err = service.fetchAndStoreDeviceGrades(deviceId, campusToken)
	if err != nil {
		log.WithError(err).Error("Could not fetch grades of device")
	}
}

func (service *Service) fetchAndStoreDeviceGrades(deviceId, campusToken string) error {
	grades, err := campus_api.FetchGrades(campusToken)
	if err != nil {
		return err
	}

	gradesRepo := ios_grades.NewRepository(service.Repository.DB)
	err = gradesRepo.EncryptAndSaveGrades(grades.Grades, deviceId, campusToken)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) fetchAndStoreDeviceExams(deviceId, campusToken string) error {
	exams, err := campus_api.FetchPersonalExams(campusToken)
	if err != nil {
		return err
	}

	examsService := ios_exams.NewService(service.Repository.DB)
	err = examsService.SaveRelevantExamsForDevice(exams.Exams, deviceId)
	if err != nil {
		return err
	}

	return nil
}

func NewService(db *Repository) *Service {
	return &Service{
		Repository: db,
	}
}
