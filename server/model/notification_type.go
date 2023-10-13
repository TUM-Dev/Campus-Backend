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
	Type         int64  `gorm:"primary_key;AUTO_INCREMENT;column:type;type:int;" json:"type"`
	Name         string `gorm:"column:name;type:text;size:65535;" json:"name"`
	Confirmation string `gorm:"column:confirmation;type:char;size:5;default:false;" json:"confirmation"`
}
