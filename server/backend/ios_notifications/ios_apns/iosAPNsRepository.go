package ios_apns

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_apns/ios_apns_jwt"
	"github.com/TUM-Dev/Campus-Backend/server/env"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	BundleId          = "de.tum.tca"
	ReadIdleTimeout   = 15 * time.Second
	HTTPClientTimeout = 60 * time.Second
)

const (
	APNsDevelopmentURL = "https://api.sandbox.push.apple.com:443"
	APNsProductionURL  = "https://api.push.apple.com:443"
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

func (r *Repository) APNsURL() string {
	if env.IsProd() {
		return APNsProductionURL
	}

	return APNsDevelopmentURL
}

func (r *Repository) CreateCampusTokenRequest(deviceId string) (*model.IOSDeviceRequestLog, error) {
	var request model.IOSDeviceRequestLog

	if err := r.DB.Create(&model.IOSDeviceRequestLog{
		DeviceID:    deviceId,
		RequestType: model.IOSBackgroundCampusTokenRequest.String(),
	}).Scan(&request).Error; err != nil {
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

	url := r.APNsURL() + "/3/device/" + notification.DeviceId

	body, _ := notification.MarshalJSON()

	client := r.httpClient

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))

	// can be e.g. alert or background
	req.Header.Set("apns-push-type", apnsPushType.String())

	req.Header.Set("apns-topic", BundleId)

	// can be a value between 1 and 10
	req.Header.Set("apns-priority", strconv.Itoa(priority))

	bearer := r.Token.GenerateNewTokenIfExpired()

	req.Header.Set("authorization", "bearer "+bearer)

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, ErrCouldNotSendNotification
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("Error while closing body: %s", err)
		}
	}(resp.Body)

	var response model.IOSRemoteNotificationResponse

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil && err != io.EOF {
		log.Error(err)
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

func NewCronRepository(db *gorm.DB) *Repository {
	token, err := ios_apns_jwt.NewToken()

	if err != nil {
		log.Fatal(err)
	}

	return NewRepository(db, token)
}
