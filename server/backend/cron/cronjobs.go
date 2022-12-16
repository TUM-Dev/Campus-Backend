package cron

import (
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type CronService struct {
	db       *gorm.DB
	gf       *gofeed.Parser
	useMensa bool
}

// names for cron jobs as specified in database
const (
	NEWS_TYPE                  = "news"
	MENSA_TYPE                 = "mensa"
	CHAT_TYPE                  = "chat"
	KINO_TYPE                  = "kino"
	ROOMFINDER_TYPE            = "roomfinder"
	TICKETSALE_TYPE            = "ticketsale"
	ALARM_TYPE                 = "alarm"
	FILE_DOWNLOAD_TYPE         = "fileDownload"
	DISH_NAME_DOWNLOAD         = "dishNameDownload"
	AVERAGE_RATING_COMPUTATION = "averageRatingComputation"
	CANTEEN_HEADCOUNT          = "canteenHeadCount"
	STORAGE_DIR                = "/Storage/" // target location of files
)

func New(db *gorm.DB, mensaCronActivated bool) *CronService {
	return &CronService{
		db:       db,
		gf:       gofeed.NewParser(),
		useMensa: mensaCronActivated,
	}
}

func (c *CronService) Run() error {
	log.Printf("running cron service. Mensa Crons Running: %t", c.useMensa)
	g := new(errgroup.Group)
	g.Go(func() error { return c.dishNameDownloadCron() })
	g.Go(func() error { return c.averageRatingComputation() })
	for {
		log.Info("Cron: checking for pending")
		var res []model.Crontab
		c.db.Model(&model.Crontab{}).
			Where("`interval` > 0 AND (lastRun+`interval`) < ? AND type IN ('news', 'fileDownload', 'averageRatingComputation','dishNameDownload', 'canteenHeadCount')", time.Now().Unix()).
			Scan(&res)

		for _, cronjob := range res {
			// Persist run to DB right away
			var offset int32 = 0
			if c.useMensa {
				if cronjob.Type.String == AVERAGE_RATING_COMPUTATION {
					if time.Now().Hour() == 16 {
						offset = 18 * 3600 // fast-forward 18 Hours to the next day + does not need to be computed overnight
					}
				}
			}

			cronjob.LastRun = int32(time.Now().Unix()) + offset
			c.db.Save(&cronjob)

			// Run each job in a separate goroutine so we can parallelize them
			switch cronjob.Type.String {
			case NEWS_TYPE:
				g.Go(func() error { return c.newsCron(&cronjob) })
			case FILE_DOWNLOAD_TYPE:
				g.Go(func() error { return c.fileDownloadCron() })
			case DISH_NAME_DOWNLOAD:
				if c.useMensa {
					g.Go(c.dishNameDownloadCron)
				}
			case AVERAGE_RATING_COMPUTATION: //call every five minutes between 11AM and 4 PM on weekdays
				if c.useMensa {
					g.Go(c.averageRatingComputation)
				}
				/*
					TODO: Implement handlers for other cronjobs
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
				*/
			case CANTEEN_HEADCOUNT:
				g.Go(func() error { return c.canteenHeadCountCron() })
			}
		}

		err := g.Wait()
		if err != nil {
			log.Println("Couldn't run all cron jobs: %v", err)
		}
		log.Info("Cron: sleeping for 60 seconds")
		time.Sleep(60 * time.Second)
	}
}
