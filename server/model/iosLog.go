package model

import "time"

const (
	IOSLogTypeGradeUpdate  = "grade_update"
	IOSLogTypeUnknown      = "unknown"
	IOSLogTypeUnknownError = "unknown_error"
	IOSLogTypeTokenRequest = "token_request"
)

// IOSLog is the model for the ios_log table.
// Should be used for logging errors and other events to the database
// to debug and analyze them later on.
type IOSLog struct {
	ID        uint `gorm:"primaryKey"`
	Data      string
	Type      string
	CreatedAt time.Time
}
