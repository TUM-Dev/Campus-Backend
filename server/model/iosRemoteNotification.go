package model

type IOSRemoteAlertNotification struct {
	Aps      IOSAlertAPS `json:"aps"`
	DeviceId string      `json:"-"`
}

type IOSAlertAPS struct {
	Alert IOSAlertAPSContent `json:"alert"`
}

type IOSAlertAPSContent struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Body     string `json:"body"`
}

type IOSRemoteNotificationResponse struct {
	Reason string `json:"reason"`
}
