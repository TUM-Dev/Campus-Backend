package subscriber

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	Repository *Repository
}

func (service *Service) NotifySubscribers(newGrades *[]model.PublishedExamResult) error {
	subscribers, err := service.Repository.FindAllSubscribers()
	if err != nil {
		return err
	}

	for _, subscriber := range *subscribers {
		if err := service.Repository.NotifySubscriber(&subscriber, newGrades); err != nil {
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
