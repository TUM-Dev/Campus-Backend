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
	NewsType                = "news"
	MensaType               = "mensa"
	ChatType                = "chat"
	KinoType                = "kino"
	RoomfinderType          = "roomfinder"
	TicketsaleType          = "ticketsale"
	AlarmType               = "alarm"
	FileDownloadType        = "fileDownload"
	LibraryReservationsType = "libraryReservations"

	StorageDir = "/Storage/" // target location of files
)

func New(db *gorm.DB) *CronService {
	return &CronService{
		db: db,
		gf: gofeed.NewParser(),
	}
}

func (c *CronService) Run() error {
	c.libraryReservationsCron()
	log.Printf("running cron service")
	for {
		log.Info("Cron: checking for pending")
		var res []model.Crontab
		c.db.Model(&model.Crontab{}).
			Where("`interval` > 0 AND (lastRun+`interval`) < ? AND type IN ('news', 'fileDownload')", time.Now().Unix()).
			Scan(&res)
		g := new(errgroup.Group)
		for _, cronjob := range res {
			// Persist run to DB right away
			cronjob.LastRun = int32(time.Now().Unix())
			c.db.Save(&cronjob)

			// Run each job in a separate goroutine so we can parallelize them
			switch cronjob.Type.String {
			case NewsType:
				g.Go(func() error { return c.newsCron(&cronjob) })
			case FileDownloadType:
				g.Go(func() error { return c.fileDownloadCron() })
			case LibraryReservationsType:
				g.Go(func() error { return c.libraryReservationsCron() })
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

