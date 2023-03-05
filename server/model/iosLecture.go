package model

import (
	"time"
)

const (
	IOSLectureSemesterWinter = "winter"
	IOSLectureSemesterSummer = "summer"
)

type IOSLecture struct {
	Id            string `gorm:"primaryKey"`
	Year          int16
	Semester      string               `gorm:"type:enum ('winter', 'summer');"`
	LastUpdate    time.Time            `gorm:"default:now()"`
	LastRequestId *string              `gorm:"default:NULL;"`
	LastRequest   *IOSDeviceRequestLog `gorm:"constraint:OnDelete:SET NULL;"`
}
