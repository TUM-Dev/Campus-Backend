package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	"github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

// NotificationConfirmation struct is a row record of the notification_confirmation table in the tca database
type NotificationConfirmation struct {
	//[ 0] notification                                   int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	Notification int32 `gorm:"primary_key;column:notification;type:int;" json:"notification"`
	//[ 1] device                                         int                  null: false  primary: true   isArray: false  auto: false  col: int             len: -1      default: []
	Device int32 `gorm:"primary_key;column:device;type:int;" json:"device"`
	//[ 2] sent                                           tinyint              null: false  primary: false  isArray: false  auto: false  col: tinyint         len: -1      default: [0]
	Sent int32 `gorm:"column:sent;type:tinyint;default:0;" json:"sent"`
	//[ 3] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 4] received                                       timestamp            null: true   primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: []
	Received null.Time `gorm:"column:received;type:timestamp;" json:"received"`
}
