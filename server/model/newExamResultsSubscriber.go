package model

import (
	"database/sql"
	"time"
)

type NewExamResultsSubscriber struct {
	CallbackUrl    string         `gorm:"primary_key" json:"callbackUrl"`
	ApiKey         sql.NullString `json:"-"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	LastNotifiedAt sql.NullTime   `json:"lastNotifiedAt"`
}
