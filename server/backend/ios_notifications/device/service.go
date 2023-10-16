// Package device provides functions to create/delete ios devices
package device

import (
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/backend/campus_api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_exams"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_grades"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Service struct {
	Repository *Repository
}

var (
	ErrCouldNotCreateDevice = status.Error(codes.Internal, "Could not create device")
	ErrCouldNotDeleteDevice = status.Error(codes.Internal, "Could not delete device")

	iosRegisteredDevices = promauto.NewGauge(prometheus.GaugeOpts{
		Subsystem: "ios",
		Name:      "ios_created_devices",
		Help:      "The number of currently created ios devices",
	})
)

func (service *Service) CreateDevice(request *pb.CreateDeviceRequest) (*pb.CreateDeviceReply, error) {
	device := model.IOSDevice{
		DeviceID:  request.GetDeviceId(),
		PublicKey: request.GetPublicKey(),
	}

	deviceAlreadyExists, err := service.Repository.CreateDevice(&device)
	if err != nil {
		return nil, ErrCouldNotCreateDevice
	}

	if !deviceAlreadyExists {
		service.handleFirstDeviceRegistration(request)
	}

	return &pb.CreateDeviceReply{
		DeviceId: request.GetDeviceId(),
	}, nil
}

func (service *Service) DeleteDevice(request *pb.DeleteDeviceRequest) (*pb.DeleteDeviceReply, error) {
	if err := service.Repository.DeleteDevice(request.GetDeviceId()); err != nil {
		return nil, ErrCouldNotDeleteDevice
	}

	iosRegisteredDevices.Dec()
	return &pb.DeleteDeviceReply{
		DeviceId: request.GetDeviceId(),
	}, nil
}

func (service *Service) handleFirstDeviceRegistration(request *pb.CreateDeviceRequest) {
	campusToken := request.GetCampusApiToken()
	deviceId := request.GetDeviceId()

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
