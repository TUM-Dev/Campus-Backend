// Package ios_device provides functions to register and remove ios devices
package ios_device

import (
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/campus_api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/influx"
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
		err := handleFirstDeviceRegistration(request)
		if err != nil {
			return nil, ErrCouldNotRegisterDevice
		}
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

func handleFirstDeviceRegistration(request *pb.RegisterDeviceRequest) error {
	influx.LogIOSRegisterDevice(request.GetDeviceId())

	campusToken := request.GetCampusApiToken()

	lectures, err := campus_api.FetchPersonalLectures(campusToken)
	if err != nil {
		return err
	}

	for _, lecture := range lectures.Lectures {
		log.Infof(lecture.LectureTitle)
	}

	grades, err := campus_api.FetchGrades(campusToken)
	if err != nil {
		return err
	}

	for _, grade := range grades.Grades {
		log.Infof(grade.LectureTitle)

	}

	return nil
}

func NewService(db *Repository) *Service {
	return &Service{
		Repository: db,
	}
}
