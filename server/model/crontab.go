package model

import (
	"github.com/guregu/null"
)

// Crontab struct is a row record of the crontab table in the tca database
type Crontab struct {
	Cron     int64       `gorm:"primary_key;autoIncrement;column:cron;type:int;" json:"cron"`
	Interval int32       `gorm:"column:interval;type:int;default:7200;not null" json:"interval"`
	LastRun  int32       `gorm:"column:lastRun;type:int;default:0;not null" json:"last_run"`
	Type     null.String `gorm:"column:type;type:enum ('news','mensa','ticketsale','alarm','fileDownload','canteenHeadCount','iosNotifications','iosActivityReset','newExamResultsHook','movie','feedbackEmail','dishNameDownload','scrapeStudentClubs');" json:"type"`
	ID       null.Int    `gorm:"column:id;type:int;" json:"id"`
}

// TableName sets the insert table name for this struct type
func (n *Crontab) TableName() string {
	return "crontab"
}
