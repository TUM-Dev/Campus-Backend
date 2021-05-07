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

// NotificationType struct is a row record of the notification_type table in the tca database
type NotificationType struct {
	//[ 0] type                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Type int32 `gorm:"primary_key;AUTO_INCREMENT;column:type;type:int;" json:"type"`
	//[ 1] name                                           text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Name string `gorm:"column:name;type:text;size:65535;" json:"name"`
	//[ 2] confirmation                                   char(5)              null: false  primary: false  isArray: false  auto: false  col: char            len: 5       default: [false]
	Confirmation string `gorm:"column:confirmation;type:char;size:5;default:false;" json:"confirmation"`
}
