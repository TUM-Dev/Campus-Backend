package model

// inspired by https://github.com/sideshow/apns2

import (
	"encoding/json"

	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/crypto"
)

type IOSNotificationPayload struct {
	content  map[string]interface{}
	DeviceId string
}

type IOSRemoteNotificationAPS struct {
	Alert            interface{} `json:"alert,omitempty"`
	Badge            interface{} `json:"badge,omitempty"`
	Category         string      `json:"category,omitempty"`
	ContentAvailable int         `json:"content-available,omitempty"`
	MutableContent   int         `json:"mutable-content,omitempty"`
	RelevanceScore   interface{} `json:"relevance-score,omitempty"`
	Sound            interface{} `json:"sound,omitempty"`
	ThreadID         string      `json:"thread-id,omitempty"`
	URLArgs          []string    `json:"url-args,omitempty"`
}

type IOSAlertAPSContent struct {
	Action          string   `json:"action,omitempty"`
	ActionLocKey    string   `json:"action-loc-key,omitempty"`
	Body            string   `json:"body,omitempty"`
	LaunchImage     string   `json:"launch-image,omitempty"`
	LocArgs         []string `json:"loc-args,omitempty"`
	LocKey          string   `json:"loc-key,omitempty"`
	Title           string   `json:"title,omitempty"`
	Subtitle        string   `json:"subtitle,omitempty"`
	TitleLocArgs    []string `json:"title-loc-args,omitempty"`
	TitleLocKey     string   `json:"title-loc-key,omitempty"`
	SummaryArg      string   `json:"summary-arg,omitempty"`
	SummaryArgCount int      `json:"summary-arg-count,omitempty"`
}

type IOSRemoteNotificationResponse struct {
	Reason string `json:"reason"`
}

type IOSAPNSPushType int64

const (
	IOSAPNSPushTypeAlert IOSAPNSPushType = iota
	IOSAPNSPushTypeBackground
)

func (pt IOSAPNSPushType) String() string {
	switch pt {
	case IOSAPNSPushTypeAlert:
		return "alert"
	case IOSAPNSPushTypeBackground:
		return "background"
	}

	return "unknown"
}

type IOSBackgroundNotificationType int64

const (
	IOSBackgroundCampusTokenRequest IOSBackgroundNotificationType = iota
)

func (requestType IOSBackgroundNotificationType) String() string {
	switch requestType {
	case IOSBackgroundCampusTokenRequest:
		return "CAMPUS_TOKEN_REQUEST"
	}

	return "unknown"
}

func NewIOSNotificationPayload(deviceId string) *IOSNotificationPayload {
	return &IOSNotificationPayload{
		content: map[string]interface{}{
			"aps": &IOSRemoteNotificationAPS{},
		},
		DeviceId: deviceId,
	}
}

func (np *IOSNotificationPayload) Alert(title string, subtitle string, body string) *IOSNotificationPayload {
	alert := np.aps().alert()

	alert.Title = title
	alert.Subtitle = subtitle
	alert.Body = body

	return np
}

func (np *IOSNotificationPayload) Background(requestId string, requestType IOSBackgroundNotificationType) *IOSNotificationPayload {
	np.aps().ContentAvailable = 1

	np.content["request_id"] = requestId
	np.content["notification_type"] = requestType.String()

	return np
}

func (np *IOSNotificationPayload) aps() *IOSRemoteNotificationAPS {
	return np.content["aps"].(*IOSRemoteNotificationAPS)
}

func (aps *IOSRemoteNotificationAPS) alert() *IOSAlertAPSContent {
	if _, ok := aps.Alert.(*IOSAlertAPSContent); !ok {
		aps.Alert = &IOSAlertAPSContent{}
	}
	return aps.Alert.(*IOSAlertAPSContent)
}

func (np *IOSNotificationPayload) Encrypt(publicKey string) *IOSNotificationPayload {
	alert := np.aps().alert()

	np.aps().MutableContent = 1

	if alert.Title != "" {
		res, err := crypto.AsymmetricEncrypt(alert.Title, publicKey)

		if err != nil {
			alert.Title = "You have a new notification"
		} else {
			alert.Title = res.String()
		}
	}

	if alert.Body != "" {
		res, err := crypto.AsymmetricEncrypt(alert.Body, publicKey)

		if err != nil {
			alert.Body = ""
		} else {
			alert.Body = res.String()
		}
	}

	if alert.Subtitle != "" {
		res, err := crypto.AsymmetricEncrypt(alert.Subtitle, publicKey)

		if err != nil {
			alert.Subtitle = ""
		} else {
			alert.Subtitle = res.String()
		}
	}

	return np
}

func (np *IOSNotificationPayload) MarshalJSON() ([]byte, error) {
	return json.Marshal(np.content)
}
