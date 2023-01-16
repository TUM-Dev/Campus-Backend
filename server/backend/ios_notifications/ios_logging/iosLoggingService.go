package ios_logging

import (
	"fmt"
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Service struct {
	Repository *Repository
}

func (service *Service) Log(data string, args ...interface{}) {
	service.LogWithType(data, model.IOSLogTypeUnknown, args...)
}

func (service *Service) LogError(data string, args ...interface{}) {
	service.LogWithType(data, model.IOSLogTypeUnknownError, args...)
}

func (service *Service) LogGradesUpdating(data string, args ...interface{}) {
	service.LogWithType(data, model.IOSLogTypeGradeUpdate, args...)
}

func (service *Service) LogTokenRequest(data string, args ...interface{}) {
	service.LogWithType(data, model.IOSLogTypeTokenRequest, args...)
}

func (service *Service) LogWithType(data string, logType string, args ...interface{}) {
	logData := model.IOSLog{
		Data: fmt.Sprintf(data, args...),
		Type: logType,
	}

	err := service.Repository.Log(&logData)

	if err != nil {
		log.Errorf("Error while logging data: %s", err)
	}
}

func NewLogger(db *gorm.DB) *Service {
	return &Service{
		Repository: NewRepository(db),
	}
}
