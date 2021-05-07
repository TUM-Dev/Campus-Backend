package cron

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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
		g := new(errgroup.Group)
		for _, cronjob := range res {
			// Persist run to DB right away
			cronjob.LastRun = int32(time.Now().Unix())
			c.db.Save(&cronjob)

			// Run each job in a separate goroutine so we can parallelize them
			switch cronjob.Type.String {
			case NEWS_TYPE:
				g.Go(func() error { return c.newsCron() })
			case MENSA_TYPE:
				g.Go(func() error { return c.mensaCron() })
			case CHAT_TYPE:
				g.Go(func() error { return c.chatCron() })
			case KINO_TYPE:
				g.Go(func() error { return c.kinoCron() })
			case ROOMFINDER_TYPE:
				g.Go(func() error { return c.roomFinderCron() })
			case TICKETSALE_TYPE:
				g.Go(func() error { return c.roomFinderCron() })
			case ALARM_TYPE:
				g.Go(func() error { return c.alarmCron() })
			}
		}
		err := g.Wait()
		if err != nil {
			log.Println("Couldn't run all cron jobs: %v", err)
		}
		log.Info("Cron: sleeping for 60 seconds")
		time.Sleep(60 * time.Second)
	}

	return nil
}
