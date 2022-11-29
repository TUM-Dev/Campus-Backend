package ios_apns

import (
	"bytes"
	"encoding/json"
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Repository struct {
	TeamId string
}

func (r *Repository) APNsURL() string {
	// production
	// return "https://api.push.apple.com:443"
	return "https://api.sandbox.push.apple.com:443"
}

func (r *Repository) SendTestNotification(notification *model.IOSRemoteAlertNotification) (*model.IOSRemoteNotificationResponse, error) {

	url := r.APNsURL() + "/3/device/" + notification.DeviceId

	body, _ := json.Marshal(notification)

	client := &http.Client{}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))

	// can be e.g. alert or background
	req.Header.Set("apns-push-type", "alert")

	req.Header.Set("apns-topic", "de.tum.tca")

	// can be a value between 1 and 10
	req.Header.Set("apns-priority", "10")

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer resp.Body.Close()

	var response model.IOSRemoteNotificationResponse

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error(err)
		return nil, err
	}

	return &response, nil
}

func NewRepository() *Repository {
	return &Repository{
		TeamId: "2J3C6P6X3N",
	}
}
