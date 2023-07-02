package model

import (
	"time"
)

type DbLecture struct {
	Id            string `gorm:"primaryKey"`
	Year          int16
	Semester      string               `gorm:"type:enum ('winter', 'summer');"`
	LastUpdate    time.Time            `gorm:"default:now()"`
	LastRequestId *string              `gorm:"default:NULL;"`
	LastRequest   *IOSDeviceRequestLog `gorm:"constraint:OnDelete:SET NULL;"`
	Title         string               `gorm:"not null"`
}

func (l *DbLecture) TableName() string {
	return "lectures"
}
