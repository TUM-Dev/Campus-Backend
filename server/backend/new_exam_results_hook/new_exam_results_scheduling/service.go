package new_exam_results_scheduling

import (
	"github.com/TUM-Dev/Campus-Backend/server/backend/campus_api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/apns"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/new_exam_results_hook/new_exam_results_subscriber"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
)

const (
	MockAPIToken = "DDF9A212B2F80A01C6D0307B8455EEAA"
)

type Service struct {
	Repository        *Repository
	DevicesRepository *device.Repository
	APNs              *apns.Service
}

func (service *Service) HandleScheduledCron() error {
	log.Info("Fetching published exam results")

	apiResult, err := campus_api.FetchExamResultsPublished(MockAPIToken)
	if err != nil {
		return err
	}

	var apiExamResults []model.ExamResultPublished
	for _, apiExamResult := range apiResult.ExamResults {
		apiExamResults = append(apiExamResults, *apiExamResult.ToDBExamResult())
	}

	storedExamResults, err := service.Repository.FindAllExamResultsPublished()
	if err != nil {
		return err
	}

	newPublishedExamResults := service.findNewPublishedExamResults(&apiExamResults, storedExamResults)

	if len(*newPublishedExamResults) > 0 {
		service.notifySubscribers(newPublishedExamResults)
	} else {
		log.Info("No new published exam results")
	}

	return service.Repository.StoreExamResultsPublished(apiExamResults)
}

func (service *Service) findNewPublishedExamResults(apiExamResults, storedExamResults *[]model.ExamResultPublished) *[]model.ExamResultPublished {
	var apiExamResultsMap = make(map[string]model.ExamResultPublished)
	for _, apiExamResult := range *apiExamResults {
		apiExamResultsMap[apiExamResult.ExamID] = apiExamResult
	}

	var storedExamResultsMap = make(map[string]model.ExamResultPublished)
	for _, storedExamResult := range *storedExamResults {
		storedExamResultsMap[storedExamResult.ExamID] = storedExamResult
	}

	var newPublishedExamResults []model.ExamResultPublished

	for id, result := range apiExamResultsMap {
		if storedResult, ok := storedExamResultsMap[id]; ok && !storedResult.Published && result.Published {
			newPublishedExamResults = append(newPublishedExamResults, result)
		}
	}

	return &newPublishedExamResults
}

func (service *Service) notifySubscribers(newPublishedExamResults *[]model.ExamResultPublished) {
	subscribersRepo := new_exam_results_subscriber.NewRepository(service.Repository.DB)
	subscribersService := new_exam_results_subscriber.NewService(subscribersRepo)

	err := subscribersService.NotifySubscribers(newPublishedExamResults)
	if err != nil {
		log.WithError(err).Error("Failed to notify subscribers")
		return
	}
}

func NewService(repository *Repository,
	devicesRepository *device.Repository,
	apnsService *apns.Service,
) *Service {
	return &Service{
		Repository:        repository,
		DevicesRepository: devicesRepository,
		APNs:              apnsService,
	}
}
