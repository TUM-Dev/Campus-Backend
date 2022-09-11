package cron

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

type CronService struct {
	db *gorm.DB
	gf *gofeed.Parser
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
	STORAGE_DIR                = "/Storage/" // target location of files
)

func New(db *gorm.DB) *CronService {
	return &CronService{
		db: db,
		gf: gofeed.NewParser(),
	}
}

func (c *CronService) Run() error {
	log.Printf("running cron service")
	g := new(errgroup.Group)
	g.Go(func() error { return c.dishNameDownloadCron() })
	g.Go(func() error { return c.averageRatingComputation() })
	for {
		log.Info("Cron: checking for pending")
		var res []model.Crontab
		c.db.Model(&model.Crontab{}).
			Where("`interval` > 0 AND (lastRun+`interval`) < ? AND type IN ('news', 'fileDownload', 'averageRatingComputation','dishNameDownload')", time.Now().Unix()).
			Scan(&res)

		for _, cronjob := range res {
			// Persist run to DB right away
			var offset int32 = 0
			if cronjob.Type.String == AVERAGE_RATING_COMPUTATION {
				if time.Now().Hour() == 16 {
					offset = 18 * 3600 // fast-forward 18 Hours to the next day + does not need to be computed overnight
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
				g.Go(func() error { return c.dishNameDownloadCron() })
			case AVERAGE_RATING_COMPUTATION: //call every five minutes between 11AM and 4 PM on weekdays
				g.Go(func() error { return c.averageRatingComputation() })
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
