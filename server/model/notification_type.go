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

// NotificationType struct is a row record of the notification_type table in the tca database
type NotificationType struct {
	Type         int64  `gorm:"primary_key;autoIncrement;column:type;type:int;" json:"type"`
	Name         string `gorm:"column:name;type:text;size:65535;" json:"name"`
	Confirmation string `gorm:"column:confirmation;type:enum('true', 'false');default:'false';" json:"confirmation"`
}

// TableName sets the insert table name for this struct type
func (n *NotificationType) TableName() string {
	return "notification_type"
}
