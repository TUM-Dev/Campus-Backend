package model

import (
	"github.com/guregu/null"
)

// Crontab struct is a row record of the crontab table in the tca database
type Crontab struct {
	Cron     int64       `gorm:"primary_key;AUTO_INCREMENT;column:cron;type:int;" json:"cron"`
	Interval int32       `gorm:"column:interval;type:int;default:7200;" json:"interval"`
	LastRun  int32       `gorm:"column:lastRun;type:int;default:0;" json:"last_run"`
	Type     null.String `gorm:"column:type;type:enum ('news', 'mensa', 'movie', 'roomfinder', 'alarm', 'fileDownload','dishNameDownload', 'iosNotifications', 'iosActivityReset', 'canteenHeadCount', 'newExamResultsHook', 'scrapeStudentClubs');" json:"type"`
	ID       null.Int    `gorm:"column:id;type:int;" json:"id"`
}
