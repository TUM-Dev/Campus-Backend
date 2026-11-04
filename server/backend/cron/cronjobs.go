package cron

import (
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type CronService struct {
	db *gorm.DB
	gf *gofeed.Parser
}

// StorageDir is the directory where files are stored
// this is a variable, so it can be changed during tests
var StorageDir = "/Storage/" // target location of files

// names for cron jobs as specified in database
const (
	NewsType         = "news"
	FileDownloadType = "fileDownload"
	DishNameDownload = "dishNameDownload"
	CanteenHeadcount = "canteenHeadCount"
	MovieType        = "movie"
	FeedbackEmail    = "feedbackEmail"
	StudentClubType  = "scrapeStudentClubs"

	/* MensaType      = "mensa"
	AlarmType      = "alarm" */
)

func New(db *gorm.DB) *CronService {
	return &CronService{
		db: db,
		gf: gofeed.NewParser(),
	}
}

func (c *CronService) Run() error {
	for {
		g := new(errgroup.Group)
		log.Trace("Cron: checking for pending")
		var res []model.Crontab

		c.db.Model(&model.Crontab{}).
			Find(&res, "`interval` > 0 AND (lastRun+`interval`) < ? AND type IN (?, ?, ?, ?, ?, ?, ?)",
				time.Now().Unix(),
				NewsType,
				FileDownloadType,
				DishNameDownload,
				CanteenHeadcount,
				MovieType,
				FeedbackEmail,
				StudentClubType,
			)

		for _, cronjob := range res {
			// Persist run to DB right away
			cronFields := log.Fields{"Cron (id)": cronjob.Cron, "type": cronjob.Type.String, "LastRun": cronjob.LastRun, "interval": cronjob.Interval, "id (not real id)": cronjob.ID.Int64}
			log.WithFields(cronFields).Trace("Running cronjob")

			cronjob.LastRun = int32(time.Now().Unix())
			c.db.Save(&cronjob)

			// Run each job in a separate goroutine, so we can parallelize them
			switch cronjob.Type.String {
			case StudentClubType:
				g.Go(func() error { return c.studentClubCron(pb.Language_German) })
				g.Go(func() error { return c.studentClubCron(pb.Language_English) })
			case NewsType:
				// if this is not copied here, this may not be threads save due to go's guarantees
				// loop variable cronjob captured by func literal (govet)
				copyCronjob := cronjob
				g.Go(func() error { return c.newsCron(&copyCronjob) })
			case FileDownloadType:
				g.Go(func() error { return c.fileDownloadCron() })
			case DishNameDownload:
				g.Go(func() error { return c.dishNameDownloadCron() })
			case MovieType:
				g.Go(func() error { return c.movieCron() })
				/*
					TODO: Implement handlers for other cronjobs
					case MensaType:
						g.Go(func() error { return c.mensaCron() })
					case KinoType:
						g.Go(func() error { return c.kinoCron() })
					case AlarmType:
						g.Go(func() error { return c.alarmCron() })
				*/
			case CanteenHeadcount:
				g.Go(func() error { return c.canteenHeadCountCron() })
			case FeedbackEmail:
				g.Go(func() error { return c.feedbackEmailCron() })
			}
		}

		if err := g.Wait(); err != nil {
			log.WithError(err).Error("Couldn't run all cron jobs")
		}
		log.Trace("Cron: sleeping for 30 seconds")
		time.Sleep(30 * time.Second)
	}
}
