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

// Notification struct is a row record of the notification table in the tca database
type Notification struct {
	Notification int64       `gorm:"primary_key;AUTO_INCREMENT;column:notification;type:int;" json:"notification"`
	Type         int32       `gorm:"column:type;type:int;" json:"type"`
	Location     null.Int    `gorm:"column:location;type:int;" json:"location"`
	Title        string      `gorm:"column:title;type:text;size:65535;" json:"title"`
	Description  string      `gorm:"column:description;type:text;size:65535;" json:"description"`
	Created      time.Time   `gorm:"column:created;type:timestamp;default:current_timestamp();" json:"created"`
	Signature    null.String `gorm:"column:signature;type:text;size:65535;" json:"signature"`
	Silent       int32       `gorm:"column:silent;type:tinyint;default:0;" json:"silent"`
}

// TableName sets the insert table name for this struct type
func (n *Notification) TableName() string {
	return "notification"
}
