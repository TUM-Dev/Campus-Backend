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

const (
	NEWS_TYPE       = "news"
	MENSA_TYPE      = "mensa"
	CHAT_TYPE       = "chat"
	KINO_TYPE       = "kino"
	ROOMFINDER_TYPE = "roomfinder"
	TICKETSALE_TYPE = "ticketsale"
	ALARM_TYPE      = "alarm"
)

func New(db *gorm.DB) *CronService {
	return &CronService{
		db: db,
	}
}

func (c *CronService) Run() error {

	for {
		log.Info("Cron: checking for pending")
		var res []model.Crontab
		c.db.Where("interval > 0 AND (lastRun+interval) < ?", time.Now().Unix()).Scan(&res)
		for _, cronjob := range res {
			switch cronjob.Type.String {
			case NEWS_TYPE:
				newsCron()
			case MENSA_TYPE:
				mensaCron()
			case CHAT_TYPE:
				chatCron()
			case KINO_TYPE:
				kinoCron()
			case ROOMFINDER_TYPE:
				roomFinderCron()
			case TICKETSALE_TYPE:
				ticketSaleCron()
			case ALARM_TYPE:
				alarmCron()
			}
			cronjob.LastRun = int32(time.Now().Unix())
			c.db.Save(&cronjob)
		}
		log.Info("Cron: sleeping for 60 seconds")
		time.Sleep(60 * time.Second)
	}

	return nil
}
