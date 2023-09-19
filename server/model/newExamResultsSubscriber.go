package model

import (
	"database/sql"
	"time"
)

type NewExamResultsSubscriber struct {
	CallbackUrl    string         `gorm:"primary_key" json:"callbackUrl"`
	ApiKey         null.String `json:"-"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	LastNotifiedAt null.Time   `json:"lastNotifiedAt"`
}
