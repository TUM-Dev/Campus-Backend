package model

import "time"

const (
	IOSLogTypeGradeUpdate  = "grade_update"
	IOSLogTypeUnknown      = "unknown"
	IOSLogTypeUnknownError = "unknown_error"
	IOSLogTypeTokenRequest = "token_request"
)

type IOSLog struct {
	ID        uint `gorm:"primaryKey"`
	Data      string
	Type      string
	CreatedAt time.Time
}
