package new_exam_results_subscriber

import (
	"bytes"
	"encoding/json"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

type Repository struct {
	DB *gorm.DB
}

func (repository *Repository) FindAllSubscribers() (*[]model.NewExamResultsSubscriber, error) {
	db := repository.DB

	var subscribers []model.NewExamResultsSubscriber

	err := db.Find(&subscribers).Error

	return &subscribers, err
}

func (repository *Repository) NotifySubscriber(subscriber *model.NewExamResultsSubscriber, newGrades *[]model.ExamResultPublished) error {
	url := subscriber.CallbackUrl

	body, err := json.Marshal(newGrades)
	if err != nil {
		log.WithError(err).Error("Error while marshalling newGrades")
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))

	req.Header.Set("Content-Type", "application/json")

	if subscriber.ApiKey.Valid {
		req.Header.Set("Authorization", subscriber.ApiKey.String)
	}

	if err != nil {
		log.WithError(err).Error("Error while creating request")
		return err
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.WithField("url", url).WithError(err).Error("Error while fetching url")
		return err
	}

	return nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
