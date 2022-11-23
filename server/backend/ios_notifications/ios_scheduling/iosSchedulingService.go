package ios_scheduling

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	Repository *Repository
}

func (service *Service) HandleScheduledCron() error {
	priorities, err := service.Repository.FindSchedulingPriorities()

	if err != nil {
		return err
	}

	currentPriority := findIOSSchedulingPriorityForNow(priorities)

	log.Info("Current priority: ", currentPriority)

	return nil
}

func findIOSSchedulingPriorityForNow(priorities []model.IOSSchedulingPriority) *model.IOSSchedulingPriority {
	var prioritiesThatAreInRange []model.IOSSchedulingPriority

	for _, priority := range priorities {
		if priority.IsCurrentlyInRange() {
			prioritiesThatAreInRange = append(prioritiesThatAreInRange, priority)
		}
	}

	if len(prioritiesThatAreInRange) == 0 {
		return model.DefaultIOSSchedulingPriority()
	}

	return mergeIOSSchedulingPriorities(prioritiesThatAreInRange)
}

func mergeIOSSchedulingPriorities(priorities []model.IOSSchedulingPriority) *model.IOSSchedulingPriority {
	mergedPriority := model.DefaultIOSSchedulingPriority()
	prioritiesSum := 0

	for _, priority := range priorities {
		if priority.IsMorePreciseThan(mergedPriority) {
			mergedPriority = &priority
		}

		prioritiesSum += priority.Priority
	}

	mergedPriority.Priority = prioritiesSum / len(priorities)

	return mergedPriority
}

func NewService(repository *Repository) *Service {
	return &Service{
		Repository: repository,
	}
}
