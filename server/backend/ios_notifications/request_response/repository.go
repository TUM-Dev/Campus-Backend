package request_response

import (
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/apns"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"gorm.io/gorm"
)

type Repository struct {
	DB    *gorm.DB
	Token *apns.JWTToken
}

func (r *Repository) SaveEncryptedGrade(grade *model.IOSEncryptedGrade) error {
	return r.DB.Create(grade).Error
}

func (r *Repository) GetIOSDeviceRequest(requestId string) (*model.IOSDeviceRequestLog, error) {
	var request model.IOSDeviceRequestLog
	if err := r.DB.First(&request, "request_id = ?", requestId).Error; err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *Repository) GetIOSEncryptedGrades(deviceId string) ([]model.IOSEncryptedGrade, error) {
	var grades []model.IOSEncryptedGrade
	if err := r.DB.Find(&grades, "device_id = ?", deviceId).Error; err != nil {
		return nil, err
	}

	return grades, nil
}

func (r *Repository) DeleteEncryptedGrades(deviceId string) error {
	return r.DB.Delete(&model.IOSEncryptedGrade{}, "device_id = ?", deviceId).Error
}

func (r *Repository) DeleteRequestLog(requestId string) error {
	return r.DB.Delete(&model.IOSDeviceRequestLog{}, "request_id = ?", requestId).Error
}

func (r *Repository) DeleteAllRequestLogsForThisDeviceWithType(requestLog *model.IOSDeviceRequestLog) error {
	return r.DB.
		Delete(
			&model.IOSDeviceRequestLog{},
			"device_id = ? and request_type = ?",
			requestLog.DeviceID,
			requestLog.RequestType,
		).Error
}

func NewRepository(db *gorm.DB, token *apns.JWTToken) *Repository {
	return &Repository{
		DB:    db,
		Token: token,
	}
}
