package cron

import (
	"net/http"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"

	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type CronService struct {
	db     *gorm.DB
	client *http.Client
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
)

func New(db *gorm.DB) *CronService {
	return &CronService{
		db:     db,
		client: webfetch.Client,
	}
}

func (c *CronService) Run() error {
	for {
		g := new(errgroup.Group)
		log.Trace("Cron: checking for pending")
		var res []model.Crontab

		if err := c.db.Model(&model.Crontab{}).
			Find(&res, "`interval` > 0 AND (lastRun+`interval`) < ? AND type IN (?, ?, ?, ?, ?, ?, ?)",
				time.Now().Unix(),
				NewsType,
				FileDownloadType,
				DishNameDownload,
				CanteenHeadcount,
				MovieType,
				FeedbackEmail,
				StudentClubType,
			).Error; err != nil {
			log.WithError(err).Error("Could not load pending cron jobs")
			time.Sleep(30 * time.Second)
			continue
		}

		for _, cronjob := range res {
			// Persist run to DB right away
			cronFields := log.Fields{"Cron (id)": cronjob.Cron, "type": cronjob.Type.String, "LastRun": cronjob.LastRun, "interval": cronjob.Interval, "id (not real id)": cronjob.ID.Int64}
			log.WithFields(cronFields).Trace("Running cronjob")

			cronjob.LastRun = int32(time.Now().Unix())
			if err := c.db.Save(&cronjob).Error; err != nil {
				log.WithFields(cronFields).WithError(err).Error("Could not persist cron run")
				continue
			}

			// Log every failure with its job, since errgroup only returns the first.
			run := func(job func() error) {
				g.Go(func() error {
					err := job()
					if err != nil {
						log.WithFields(cronFields).WithError(err).Error("Cron job failed")
					}
					return err
				})
			}
			// Run each job in a separate goroutine.
			switch cronjob.Type.String {
			case StudentClubType:
				run(func() error { return c.studentClubCron(pb.Language_German) })
				run(func() error { return c.studentClubCron(pb.Language_English) })
			case NewsType:
				// if this is not copied here, this may not be threads save due to go's guarantees
				// loop variable cronjob captured by func literal (govet)
				copyCronjob := cronjob
				run(func() error { return c.newsCron(&copyCronjob) })
			case FileDownloadType:
				run(func() error { return c.fileDownloadCron() })
			case DishNameDownload:
				run(func() error { return c.dishNameDownloadCron() })
			case MovieType:
				run(func() error { return c.movieCron() })
			case CanteenHeadcount:
				run(func() error { return c.canteenHeadCountCron() })
			case FeedbackEmail:
				run(func() error { return c.feedbackEmailCron() })
			}
		}

		if err := g.Wait(); err != nil {
			log.WithError(err).Error("Couldn't run all cron jobs")
		}
		log.Trace("Cron: sleeping for 30 seconds")
		time.Sleep(30 * time.Second)
	}
}
