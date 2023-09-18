package model

import (
	"time"

	"github.com/guregu/null"
)

// Devices struct is a row record of the devices table in the tca database
type Devices struct {
	Device          int32       `gorm:"primary_key;AUTO_INCREMENT;column:device;type:int;" json:"device"`
	Member          null.Int    `gorm:"column:member;type:int;" json:"member"`
	UUID            string      `gorm:"column:uuid;type:varchar(50);" json:"uuid"`
	Created         null.Time   `gorm:"column:created;type:timestamp;" json:"created"`
	LastAccess      time.Time   `gorm:"column:lastAccess;type:timestamp;default:0000-00-00 00:00:00;" json:"last_access"`
	LastAPI         string      `gorm:"column:lastApi;type:text;size:16777215;" json:"last_api"`
	Developer       string      `gorm:"column:developer;type:char;size:5;default:false;" json:"developer"`
	OsVersion       string      `gorm:"column:osVersion;type:text;size:16777215;" json:"os_version"`
	AppVersion      string      `gorm:"column:appVersion;type:text;size:16777215;" json:"app_version"`
	Counter         int32       `gorm:"column:counter;type:int;default:0;" json:"counter"`
	Pk              null.String `gorm:"column:pk;type:text;size:4294967295;" json:"pk"`
	PkActive        string      `gorm:"column:pkActive;type:char;size:5;default:false;" json:"pk_active"`
	GcmToken        null.String `gorm:"column:gcmToken;type:text;size:65535;" json:"gcm_token"`
	GcmStatus       null.String `gorm:"column:gcmStatus;type:varchar(200);" json:"gcm_status"`
	ConfirmationKey null.String `gorm:"column:confirmationKey;type:varchar(35);" json:"confirmation_key"`
	KeyCreated      null.Time   `gorm:"column:keyCreated;type:datetime;" json:"key_created"`
	KeyConfirmed    null.Time   `gorm:"column:keyConfirmed;type:datetime;" json:"key_confirmed"`
}
