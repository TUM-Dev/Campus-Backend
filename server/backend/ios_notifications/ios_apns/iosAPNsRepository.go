package ios_apns

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns/ios_apns_jwt"
	"github.com/TUM-Dev/Campus-Backend/server/env"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"gorm.io/gorm"
)

const (
	// BundleId from the Apple Developer Portal
	BundleId = "de.tum.tca"
	// ReadIdleTimeout is the idle time after which the http2 transport will do a health check
	ReadIdleTimeout = 15 * time.Second
	// HTTPClientTimeout is the timeout for the http client used to send notifications
	HTTPClientTimeout = 60 * time.Second
)

const (
	ApnsDevelopmentURL = "https://api.sandbox.push.apple.com:443"
	ApnsProductionURL  = "https://api.push.apple.com:443"
)

var (
	ErrCouldNotSendNotification   = errors.New("could not send notification")
	ErrCouldNotDecodeAPNsResponse = errors.New("could not decode apns response")
)

type Repository struct {
	DB         gorm.DB
	Token      *ios_apns_jwt.Token
	httpClient *http.Client
}

// ApnsUrl uses the environment variable ENVIRONMENT to determine whether
// to use the production or development APNs URL.
func (r *Repository) ApnsUrl() string {
	if env.IsProd() {
		return ApnsProductionURL
	}
	return ApnsDevelopmentURL
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

	url := r.ApnsUrl() + "/3/device/" + notification.DeviceId
	body, _ := notification.MarshalJSON()

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))

	// can be e.g. alert or background
	req.Header.Set("apns-push-type", apnsPushType.String())
	req.Header.Set("apns-topic", BundleId)
	// can be a value between 1 and 10
	req.Header.Set("apns-priority", strconv.Itoa(priority))

	bearer := r.Token.GenerateNewTokenIfExpired()
	req.Header.Set("authorization", "bearer "+bearer)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		log.WithError(err).Error("Could not send notification")
		return nil, ErrCouldNotSendNotification
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.WithError(err).Error("Could not close body")
		}
	}(resp.Body)

	var response model.IOSRemoteNotificationResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil && err != io.EOF {
		log.WithError(err).Error("Could not decode APNs response")
		return nil, ErrCouldNotDecodeAPNsResponse
	}

	return &response, nil
}

func NewRepository(db *gorm.DB, token *ios_apns_jwt.Token) *Repository {
	transport := &http2.Transport{
		ReadIdleTimeout: ReadIdleTimeout,
	}

	return &Repository{
		DB:    *db,
		Token: token,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   HTTPClientTimeout,
		},
	}
}

func NewCronRepository(db *gorm.DB) (*Repository, error) {
	if err := ValidateRequirementsForIOSNotificationsService(); err != nil {
		log.WithError(err).Warn("Failed to validate requirements for ios notifications service")
		return nil, err
	}

	token, err := ios_apns_jwt.NewToken()
	if err != nil {
		log.WithError(err).Error("Could not create APNs token")
	}

	return NewRepository(db, token), nil
}
