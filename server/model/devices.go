package model

import (
	"time"

	"github.com/guregu/null"
)

// Devices struct is a row record of the devices table in the tca database
type Devices struct {
	//[ 0] device                                         int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Device int32 `gorm:"primary_key;AUTO_INCREMENT;column:device;type:int;" json:"device"`
	//[ 1] member                                         int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Member null.Int `gorm:"column:member;type:int;" json:"member"`
	//[ 2] uuid                                           varchar(50)          null: false  primary: false  isArray: false  auto: false  col: varchar         len: 50      default: []
	UUID string `gorm:"column:uuid;type:varchar(50);" json:"uuid"`
	//[ 3] created                                        timestamp            null: true   primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: []
	Created null.Time `gorm:"column:created;type:timestamp;" json:"created"`
	//[ 4] lastAccess                                     timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [0000-00-00 00:00:00]
	LastAccess time.Time `gorm:"column:lastAccess;type:timestamp;default:0000-00-00 00:00:00;" json:"last_access"`
	//[ 5] lastApi                                        text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	LastAPI string `gorm:"column:lastApi;type:text;size:16777215;" json:"last_api"`
	//[ 6] developer                                      char(5)              null: false  primary: false  isArray: false  auto: false  col: char            len: 5       default: [false]
	Developer string `gorm:"column:developer;type:char;size:5;default:false;" json:"developer"`
	//[ 7] osVersion                                      text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	OsVersion string `gorm:"column:osVersion;type:text;size:16777215;" json:"os_version"`
	//[ 8] appVersion                                     text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	AppVersion string `gorm:"column:appVersion;type:text;size:16777215;" json:"app_version"`
	//[ 9] counter                                        int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	Counter int32 `gorm:"column:counter;type:int;default:0;" json:"counter"`
	//[10] pk                                             text(4294967295)     null: true   primary: false  isArray: false  auto: false  col: text            len: 4294967295 default: []
	Pk null.String `gorm:"column:pk;type:text;size:4294967295;" json:"pk"`
	//[11] pkActive                                       char(5)              null: false  primary: false  isArray: false  auto: false  col: char            len: 5       default: [false]
	PkActive string `gorm:"column:pkActive;type:char;size:5;default:false;" json:"pk_active"`
	//[12] gcmToken                                       text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	GcmToken null.String `gorm:"column:gcmToken;type:text;size:65535;" json:"gcm_token"`
	//[13] gcmStatus                                      varchar(200)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 200     default: []
	GcmStatus null.String `gorm:"column:gcmStatus;type:varchar(200);" json:"gcm_status"`
	//[14] confirmationKey                                varchar(35)          null: true   primary: false  isArray: false  auto: false  col: varchar         len: 35      default: []
	ConfirmationKey null.String `gorm:"column:confirmationKey;type:varchar(35);" json:"confirmation_key"`
	//[15] keyCreated                                     datetime             null: true   primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	KeyCreated null.Time `gorm:"column:keyCreated;type:datetime;" json:"key_created"`
	//[16] keyConfirmed                                   datetime             null: true   primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	KeyConfirmed null.Time `gorm:"column:keyConfirmed;type:datetime;" json:"key_confirmed"`
}
