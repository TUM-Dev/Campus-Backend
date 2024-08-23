package model

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

// NotificationConfirmation struct is a row record of the notification_confirmation table in the tca database
type NotificationConfirmation struct {
	Notification int64     `gorm:"primary_key;column:notification;type:int;" json:"notification"`
	Device       int64     `gorm:"primary_key;column:device;type:int;" json:"device"`
	Sent         int32     `gorm:"column:sent;type:tinyint;default:0;" json:"sent"`
	Created      time.Time `gorm:"column:created;type:timestamp;default:current_timestamp();" json:"created"`
	Received     null.Time `gorm:"column:received;type:timestamp;" json:"received"`
}

// TableName sets the insert table name for this struct type
func (n *NotificationConfirmation) TableName() string {
	return "notification_confirmation"
}
