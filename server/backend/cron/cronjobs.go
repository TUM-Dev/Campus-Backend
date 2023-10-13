package cron

import (
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/apns"
	"github.com/TUM-Dev/Campus-Backend/server/env"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type CronService struct {
	db   *gorm.DB
	gf   *gofeed.Parser
	APNs *apns.Service
}

// StorageDir is the directory where files are stored
// this is a variable, so it can be changed during tests
var StorageDir = "/Storage/" // target location of files

// names for cron jobs as specified in database
const (
	NewsType                 = "news"
	FileDownloadType         = "fileDownload"
	DishNameDownload         = "dishNameDownload"
	AverageRatingComputation = "averageRatingComputation"
	CanteenHeadcount         = "canteenHeadCount"
	NewExamResultsHook       = "newExamResultsHook"
	MovieType                = "movie"
	FeedbackEmail            = "feedbackEmail"

	/* MensaType      = "mensa"
	AlarmType      = "alarm" */
)

func New(db *gorm.DB) *CronService {
	return &CronService{
		db:   db,
		gf:   gofeed.NewParser(),
		APNs: apns.NewCronService(db),
	}
}

func (c *CronService) Run() error {
	log.WithField("MensaCronActive", env.IsMensaCronActive()).Debug("running cron service")
	g := new(errgroup.Group)

	if env.IsMensaCronActive() {
		g.Go(func() error { return c.dishNameDownloadCron() })
		g.Go(func() error { return c.averageRatingComputation() })
	}

	for {
		log.Trace("Cron: checking for pending")
		var res []model.Crontab

		c.db.Model(&model.Crontab{}).
			Where("`interval` > 0 AND (lastRun+`interval`) < ? AND type IN (?, ?, ?, ?, ?, ?, ?, ?)",
				time.Now().Unix(),
				NewsType,
				FileDownloadType,
				AverageRatingComputation,
				DishNameDownload,
				CanteenHeadcount,
				NewExamResultsHook,
				MovieType,
				FeedbackEmail,
			).
			Scan(&res)

		for _, cronjob := range res {
			// Persist run to DB right away
			var offset int32 = 0
			if env.IsMensaCronActive() {
				if cronjob.Type.String == AverageRatingComputation {
					if time.Now().Hour() == 16 {
						offset = 18 * 3600 // fast-forward 18 Hours to the next day + does not need to be computed overnight
					}
				}
			}
			cronFields := log.Fields{"Cron (id)": cronjob.Cron, "type": cronjob.Type.String, "offset": offset, "LastRun": cronjob.LastRun, "interval": cronjob.Interval, "id (not real id)": cronjob.ID.Int64}
			log.WithFields(cronFields).Trace("Running cronjob")

			cronjob.LastRun = int32(time.Now().Unix()) + offset
			c.db.Save(&cronjob)

			// Run each job in a separate goroutine, so we can parallelize them
			switch cronjob.Type.String {
			case NewsType:
				// if this is not copied here, this may not be threads save due to go's guarantees
				// loop variable cronjob captured by func literal (govet)
				copyCronjob := cronjob
				g.Go(func() error { return c.newsCron(&copyCronjob) })
			case FileDownloadType:
				g.Go(func() error { return c.fileDownloadCron() })
			case DishNameDownload:
				if env.IsMensaCronActive() {
					g.Go(c.dishNameDownloadCron)
				}
			case AverageRatingComputation: //call every five minutes between 11AM and 4 PM on weekdays
				if env.IsMensaCronActive() {
					g.Go(c.averageRatingComputation)
				}
			case NewExamResultsHook:
				g.Go(func() error { return c.newExamResultsHookCron() })
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
		log.Trace("Cron: sleeping for 60 seconds")
		time.Sleep(60 * time.Second)
	}
}
