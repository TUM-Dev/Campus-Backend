package model

import "time"

const (
	ServiceStatusOK          = "ok"
	ServiceStatusDegraded    = "degraded"
	ServiceStatusMaintenance = "maintenance"
)

// ServiceStatusState persists the last observed TUM IT service status to detect outage/recovery transitions.
type ServiceStatusState struct {
	ServiceKey string    `gorm:"primaryKey;column:service_key;type:varchar(32);not null;" json:"service_key"`
	Status     string    `gorm:"column:status;type:enum('ok','degraded','maintenance');not null;default:'ok';" json:"status"`
	MessageDE  string    `gorm:"column:message_de;type:varchar(512);not null;default:'';" json:"message_de"`
	MessageEN  string    `gorm:"column:message_en;type:varchar(512);not null;default:'';" json:"message_en"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp;not null;default:current_timestamp() ON UPDATE current_timestamp();" json:"updated_at"`
}
