package model

import (
	"time"

	"github.com/guregu/null"
)

// Device struct is a row record of the devices table in the tca database
type Device struct {
	Device          int64       `gorm:"primary_key;AUTO_INCREMENT;column:device;type:int;" json:"device"`
	Member          null.Int    `gorm:"column:member;type:int;" json:"member"`
	UUID            string      `gorm:"column:uuid;type:varchar(50);not null" json:"uuid"`
	Created         null.Time   `gorm:"column:created;type:timestamp;" json:"created"`
	LastAccess      time.Time   `gorm:"column:lastAccess;type:timestamp;default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP;not null" json:"last_access"`
	LastAPI         string      `gorm:"column:lastApi;type:text;not null;default:('')" json:"last_api"`
	Developer       string      `gorm:"column:developer;type:enum('true','false');default:'false';not null" json:"developer"`
	OsVersion       string      `gorm:"column:osVersion;type:text;not null;default:('')" json:"os_version"`
	AppVersion      string      `gorm:"column:appVersion;type:text;not null;default:('')" json:"app_version"`
	Counter         int32       `gorm:"column:counter;type:int;default:0;not null" json:"counter"`
	Pk              null.String `gorm:"column:pk;type:text;" json:"pk"`
	PkActive        string      `gorm:"column:pkActive;type:enum('true', 'false');default:'false';not null" json:"pk_active"`
	GcmToken        null.String `gorm:"column:gcmToken;type:text;size:65535;" json:"gcm_token"`
	GcmStatus       null.String `gorm:"column:gcmStatus;type:varchar(200);" json:"gcm_status"`
	ConfirmationKey null.String `gorm:"column:confirmationKey;type:varchar(35);" json:"confirmation_key"`
	KeyCreated      null.Time   `gorm:"column:keyCreated;type:datetime;" json:"key_created"`
	KeyConfirmed    null.Time   `gorm:"column:keyConfirmed;type:datetime;" json:"key_confirmed"`
}
