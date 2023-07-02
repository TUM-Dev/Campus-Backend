package ios_new_exams_callback

import (
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_exams"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	Repository *Repository
	APNs       *ios_apns.Service
	isActive   bool
}

func (service *Service) HandleNewExamsCallback(req *pb.NewExamsHookRequest) error {
	publishedExams := req.GetPublishedExams()

	if !service.isActive {
		log.Warn("iOS notifications are not active")
		return nil
	}

	if len(publishedExams) == 0 {
		log.Info("No exams were published")
		return nil
	}

	examIds := make([]string, len(publishedExams))

	for _, exam := range publishedExams {
		examIds = append(examIds, exam.ExamId)
	}

	examRepository := ios_exams.NewRepository(service.Repository.DB)

	devices, err := examRepository.GetDevicesThatHaveExams(&examIds)
	if err != nil {
		log.WithError(err).Info("Couldn't query all devices which wrote exams")
		return err
	}

	log.Infof("Found %d devices which wrote exams", len(*devices))

	for _, device := range *devices {
		err := service.APNs.RequestGradeUpdateForDevice(device.DeviceId)
		if err != nil {
			log.WithError(err).Infof("Couldn't request grade update for device %s", device.DeviceId)
			continue
		}
	}

	return nil
}

func NewService(repo *Repository, apns *ios_apns.Service, isActive bool) *Service {
	return &Service{
		Repository: repo,
		APNs:       apns,
		isActive:   isActive,
	}
}
