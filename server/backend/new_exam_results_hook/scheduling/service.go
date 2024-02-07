package scheduling

import (
	"github.com/TUM-Dev/Campus-Backend/server/backend/campus_api"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/apns"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/new_exam_results_hook/subscriber"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	Repository        *Repository
	DevicesRepository *device.Repository
	APNs              *apns.Service
}

func (service *Service) HandleScheduledCron() error {
	log.Info("Fetching published exam results")

	apiResult, err := campus_api.FetchExamResultsPublished()
	if err != nil {
		return err
	}

	var apiExamResults []model.PublishedExamResult
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

func (service *Service) findNewPublishedExamResults(apiExamResults, storedExamResults *[]model.PublishedExamResult) *[]model.PublishedExamResult {
	var apiExamResultsMap = make(map[string]model.PublishedExamResult)
	for _, apiExamResult := range *apiExamResults {
		apiExamResultsMap[apiExamResult.ExamID] = apiExamResult
	}

	var storedExamResultsMap = make(map[string]model.PublishedExamResult)
	for _, storedExamResult := range *storedExamResults {
		storedExamResultsMap[storedExamResult.ExamID] = storedExamResult
	}

	var newPublishedExamResults []model.PublishedExamResult

	for id, result := range apiExamResultsMap {
		if storedResult, ok := storedExamResultsMap[id]; ok && !storedResult.Published && result.Published {
			newPublishedExamResults = append(newPublishedExamResults, result)
		}
	}

	return &newPublishedExamResults
}

func (service *Service) notifySubscribers(newPublishedExamResults *[]model.PublishedExamResult) {
	log.Infof("Notifying subscribers about %d published exam results", len(*newPublishedExamResults))

	subscribersRepo := subscriber.NewRepository(service.Repository.DB)
	subscribersService := subscriber.NewService(subscribersRepo)

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