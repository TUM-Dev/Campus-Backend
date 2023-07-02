package model

import (
	"database/sql"
	"time"
)

type NewExamResultsSubscriber struct {
	ID             int32          `gorm:"primary_key;AUTO_INCREMENT;" json:"id"`
	CallbackUrl    string         `json:"callbackUrl"`
	ApiKey         sql.NullString `json:"-"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	LastNotifiedAt sql.NullTime   `json:"lastNotifiedAt"`
}
