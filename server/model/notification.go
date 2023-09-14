package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

// Notification struct is a row record of the notification table in the tca database
type Notification struct {
	//[ 0] notification                                   int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Notification int32 `gorm:"primary_key;AUTO_INCREMENT;column:notification;type:int;" json:"notification"`
	//[ 1] type                                           int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Type int32 `gorm:"column:type;type:int;" json:"type"`
	//[ 2] location                                       int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Location null.Int `gorm:"column:location;type:int;" json:"location"`
	//[ 3] title                                          text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Title string `gorm:"column:title;type:text;size:65535;" json:"title"`
	//[ 4] description                                    text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Description string `gorm:"column:description;type:text;size:65535;" json:"description"`
	//[ 5] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 6] signature                                      text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Signature null.String `gorm:"column:signature;type:text;size:65535;" json:"signature"`
	//[ 7] silent                                         tinyint              null: false  primary: false  isArray: false  auto: false  col: tinyint         len: -1      default: [0]
	Silent int32 `gorm:"column:silent;type:tinyint;default:0;" json:"silent"`
}
