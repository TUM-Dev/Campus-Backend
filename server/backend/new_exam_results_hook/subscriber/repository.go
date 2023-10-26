package subscriber

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (repository *Repository) FindAllSubscribers() (*[]model.NewExamResultsSubscriber, error) {
	var subscribers []model.NewExamResultsSubscriber
	err := repository.DB.Find(&subscribers).Error

	return &subscribers, err
}

func (repository *Repository) NotifySubscriber(subscriber *model.NewExamResultsSubscriber, newGrades *[]model.PublishedExamResult) error {
	body, err := json.Marshal(newGrades)
	if err != nil {
		log.WithError(err).Error("Error while marshalling newGrades")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, subscriber.CallbackUrl, bytes.NewBuffer(body))
	if err != nil {
		log.WithError(err).Error("Error while creating request")
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if subscriber.ApiKey.Valid {
		req.Header.Set("Authorization", subscriber.ApiKey.String)
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.WithField("url", subscriber.CallbackUrl).WithError(err).Error("Error while fetching url")
		return err
	}

	return nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
