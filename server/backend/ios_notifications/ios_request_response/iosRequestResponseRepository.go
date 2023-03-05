package ios_request_response

import (
	"database/sql"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns/ios_apns_jwt"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	DB    *gorm.DB
	Token *ios_apns_jwt.Token
}

func (r *Repository) GetIOSDeviceRequest(requestId string) (*model.IOSDeviceRequestLog, error) {
	var request model.IOSDeviceRequestLog
	if err := r.DB.First(&request, "request_id = ?", requestId).Error; err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *Repository) DeleteRequestLog(requestId string) error {
	if err := r.DB.Delete(&model.IOSDeviceRequestLog{}, "request_id = ?", requestId).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteAllRequestLogsForThisDeviceWithType(requestLog *model.IOSDeviceRequestLog) error {
	res := r.DB.
		Delete(
			&model.IOSDeviceRequestLog{},
			"device_id = ? and request_type = ?",
			requestLog.DeviceID,
			requestLog.RequestType,
		)

	if err := res.Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) SetRequestLogAsHandled(requestLog *model.IOSDeviceRequestLog) error {
	requestLog.HandledAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}

	if err := r.DB.Save(&requestLog).Error; err != nil {
		return err
	}

	return nil
}

func NewRepository(db *gorm.DB, token *ios_apns_jwt.Token) *Repository {
	return &Repository{
		DB:    db,
		Token: token,
	}
}
