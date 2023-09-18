package model

import (
	"github.com/guregu/null"
)

// TableName overrides the table name used by Crontab to `crontab` (Would otherwise auto-migrate to crontabs)
func (Crontab) TableName() string {
	return "crontab"
}

// Crontab struct is a row record of the crontab table in the tca database
type Crontab struct {
	Cron     int32       `gorm:"primary_key;AUTO_INCREMENT;column:cron;type:int;" json:"cron"`
	Interval int32       `gorm:"column:interval;type:int;default:7200;" json:"interval"`
	LastRun  int32       `gorm:"column:lastRun;type:int;default:0;" json:"last_run"`
	Type     null.String `gorm:"column:type;type:enum ('news', 'mensa', 'movie', 'roomfinder', 'alarm', 'fileDownload','dishNameDownload','averageRatingComputation', 'iosNotifications', 'iosActivityReset', 'canteenHeadCount');" json:"type"`
	ID       null.Int    `gorm:"column:id;type:int;" json:"id"`
}
