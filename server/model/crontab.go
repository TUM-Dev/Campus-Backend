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
	//[ 0] cron                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Cron int32 `gorm:"primary_key;AUTO_INCREMENT;column:cron;type:int;" json:"cron"`
	//[ 1] interval                                       int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [7200]
	Interval int32 `gorm:"column:interval;type:int;default:7200;" json:"interval"`
	//[ 2] lastRun                                        int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	LastRun int32 `gorm:"column:lastRun;type:int;default:0;" json:"last_run"`
	//[ 3] type                                           char(10)             null: true   primary: false  isArray: false  auto: false  col: char            len: 10      default: []
	Type null.String `gorm:"column:type;type:enum ('news', 'mensa', 'chat', 'kino', 'roomfinder', 'ticketsale', 'alarm', 'fileDownload','dishNameDownload','averageRatingComputation', 'canteenHeadCount', 'newExamResultsHook');" json:"type"`
	//[ 4] id                                             int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	ID null.Int `gorm:"column:id;type:int;" json:"id"`
}
