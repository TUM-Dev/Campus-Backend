package apns

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/env"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"gorm.io/gorm"
)

type Repository struct {
	DB         gorm.DB
	Token      *JWTToken
	httpClient *http.Client
}

// ApnsUrl uses the environment variable ENVIRONMENT to determine whether
// to use the production or development APNs URL.
func (r *Repository) ApnsUrl(DeviceId string) string {
	if env.IsProd() {
		return "https://api.push.apple.com:443/3/device/" + DeviceId
	}
	return "https://api.sandbox.push.apple.com:443/3/device/" + DeviceId
}

// CreateCampusTokenRequest creates a request log in the database that can be referred to
// when the app responds to the background notification.
func (r *Repository) CreateCampusTokenRequest(deviceId string) (*model.IOSDeviceRequestLog, error) {
	return r.CreateRequest(deviceId, model.IOSBackgroundCampusTokenRequest)
}

func (r *Repository) CreateRequest(deviceId string, requestType model.IOSBackgroundNotificationType) (*model.IOSDeviceRequestLog, error) {
	var request model.IOSDeviceRequestLog

	tx := r.DB.Raw(`
		insert into ios_device_request_logs (device_id, request_type)
		values (?, ?)
		returning device_id, request_id, request_type;
	`, deviceId, requestType.String()).Scan(&request)

	if err := tx.Error; err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *Repository) SendAlertNotification(payload *model.IOSNotificationPayload) (*model.IOSRemoteNotificationResponse, error) {
	return r.SendNotification(payload, model.IOSAPNSPushTypeAlert, 10)
}

func (r *Repository) SendBackgroundNotification(payload *model.IOSNotificationPayload) (*model.IOSRemoteNotificationResponse, error) {
	return r.SendNotification(payload, model.IOSAPNSPushTypeBackground, 10)
}

func (r *Repository) SendNotification(notification *model.IOSNotificationPayload, apnsPushType model.IOSAPNSPushType, priority int) (*model.IOSRemoteNotificationResponse, error) {
	body, _ := notification.MarshalJSON()

	req, _ := http.NewRequest(http.MethodPost, r.ApnsUrl(notification.DeviceId), bytes.NewBuffer(body))

	// can be e.g. alert or background
	req.Header.Set("apns-push-type", apnsPushType.String())
	req.Header.Set("apns-topic", "de.tum.tca")
	// can be a value between 1 and 10
	req.Header.Set("apns-priority", strconv.Itoa(priority))

	bearer := r.Token.GenerateNewTokenIfExpired()
	req.Header.Set("authorization", "bearer "+bearer)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("Could not send notification")
		return nil, errors.New("could not send notification")
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.WithError(err).Error("Could not close body")
		}
	}(resp.Body)

	var response model.IOSRemoteNotificationResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil && err != io.EOF {
		log.WithError(err).Error("Could not decode APNs response")
		return nil, errors.New("could not decode apns response")
	}

	return &response, nil
}

func NewRepository(db *gorm.DB, token *JWTToken) *Repository {
	transport := &http2.Transport{
		ReadIdleTimeout: 15 * time.Second,
	}

	return &Repository{
		DB:    *db,
		Token: token,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   60 * time.Second,
		},
	}
}

func NewCronRepository(db *gorm.DB) (*Repository, error) {
	if err := ValidateRequirementsForIOSNotificationsService(); err != nil {
		log.WithError(err).Warn("Failed to validate requirements for ios notifications service")
		return nil, err
	}

	token, err := NewToken()
	if err != nil {
		log.WithError(err).Error("Could not create APNs token")
	}

	return NewRepository(db, token), nil
}
