package cron

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type CronService struct {
	db *gorm.DB
}

func New(db *gorm.DB) *CronService {
	return &CronService{
		db: db,
	}
}

func (c *CronService) Run() error {

	for {
		log.Info("Cron: checking for pending")
		var res *model.Crontab
		c.db.Where("interval>0 AND (lastRun+interval) < ?", "jinzhu").Scan(res)

		log.Info("Cron: sleeping for 60 seconds")
		time.Sleep(60 * time.Second)
	}

	return nil
}
