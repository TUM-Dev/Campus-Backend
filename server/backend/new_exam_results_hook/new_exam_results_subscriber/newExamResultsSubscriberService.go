package new_exam_results_subscriber

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	Repository *Repository
}

func (service *Service) NotifySubscribers(newGrades *[]model.PublishedExamResult) error {
	repository := service.Repository

	subscribers, err := repository.FindAllSubscribers()
	if err != nil {
		return err
	}

	for _, subscriber := range *subscribers {
		err := repository.NotifySubscriber(&subscriber, newGrades)
		if err != nil {
			log.WithError(err).Error("Failed to notify subscriber")
			continue
		}
	}

	return nil
}

func NewService(repository *Repository) *Service {
	return &Service{
		Repository: repository,
	}
}
