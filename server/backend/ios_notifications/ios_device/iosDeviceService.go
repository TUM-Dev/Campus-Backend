// Package ios_device provides functions to register and remove ios devices
package ios_device

import (
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/campus_api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/influx"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_grades"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_lectures"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
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
		CreatedAt: time.Now(),
	}

	deviceAlreadyExisted, err := service.Repository.RegisterDevice(&device)

	if err != nil {
		return nil, ErrCouldNotRegisterDevice
	}

	if !deviceAlreadyExisted {
		service.handleFirstDeviceRegistration(request)
	}

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

	err := service.fetchDeviceLectures(deviceId, campusToken)
	if err != nil {
		log.WithError(err).Error("Could not fetch lectures of device")
	}

	err = service.fetchDeviceGrades(deviceId, campusToken)
	if err != nil {
		log.WithError(err).Error("Could not fetch grades of device")
	}
}

func (service *Service) fetchDeviceGrades(deviceId, campusToken string) error {
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

func (service *Service) fetchDeviceLectures(deviceId, campusToken string) error {
	lectures, err := campus_api.FetchPersonalLectures(campusToken)
	if err != nil {
		return err
	}

	lecturesService := ios_lectures.NewService(service.Repository.DB)
	err = lecturesService.SaveRelevantLecturesForDevice(lectures.Lectures, deviceId)
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
