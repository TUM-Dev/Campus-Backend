package cron

import (
	"errors"
	"net/http"
	"sync"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"

	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type CronService struct {
	db         *gorm.DB
	client     *http.Client
	newsClient *http.Client
	mu         sync.Mutex
	status     Status
}

type JobStatus struct {
	StartedAt   time.Time
	FinishedAt  time.Time
	LastSuccess time.Time
	State       string
}

// Status contains observations since this process started, not historical job results.
type Status struct {
	LastPoll time.Time
	PollOK   bool
	Jobs     map[int64]JobStatus
}

func (c *CronService) Status() Status {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := Status{LastPoll: c.status.LastPoll, PollOK: c.status.PollOK, Jobs: make(map[int64]JobStatus)}
	for id, job := range c.status.Jobs {
		result.Jobs[id] = job
	}
	return result
}

func (c *CronService) recordPoll(ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status.LastPoll, c.status.PollOK = time.Now(), ok
}

func (c *CronService) recordJob(id int64, finished bool, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.status.Jobs == nil {
		c.status.Jobs = make(map[int64]JobStatus)
	}
	job := c.status.Jobs[id]
	if !finished {
		job.StartedAt, job.FinishedAt, job.State = time.Now(), time.Time{}, "Running"
	} else {
		job.FinishedAt, job.State = time.Now(), "Succeeded"
		if err != nil {
			job.State = "Failed"
		} else {
			job.LastSuccess = job.FinishedAt
		}
	}
	c.status.Jobs[id] = job
}

func (c *CronService) claimJob(job *model.Crontab) (bool, error) {
	started := int32(time.Now().Unix())
	result := c.db.Model(job).Where("`interval` = ? AND lastRun = ?", job.Interval, job.LastRun).
		UpdateColumn("lastRun", started)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	job.LastRun = started
	return true, nil
}

// StorageDir is the directory where files are stored
// this is a variable, so it can be changed during tests
var StorageDir = "/Storage/" // target location of files

// names for cron jobs as specified in database
const (
	NewsType            = "news"
	FileDownloadType    = "fileDownload"
	DishNameDownload    = "dishNameDownload"
	CanteenHeadcount    = "canteenHeadCount"
	MovieType           = "movie"
	FeedbackEmail       = "feedbackEmail"
	StudentClubType     = "scrapeStudentClubs"
	TUMOnlineStatusType = "tumonlineStatus"
)

func New(db *gorm.DB) *CronService {
	return &CronService{
		db:         db,
		client:     webfetch.Client,
		newsClient: webfetch.PublicClient,
	}
}

func (c *CronService) Run() error {
	for {
		g := new(errgroup.Group)
		log.Trace("Cron: checking for pending")
		var res []model.Crontab

		if err := c.db.Model(&model.Crontab{}).
			Find(&res, "`interval` > 0 AND (lastRun+`interval`) < ? AND type IN (?, ?, ?, ?, ?, ?, ?, ?)",
				time.Now().Unix(),
				NewsType,
				FileDownloadType,
				DishNameDownload,
				CanteenHeadcount,
				MovieType,
				FeedbackEmail,
				StudentClubType,
				TUMOnlineStatusType,
			).Error; err != nil {
			c.recordPoll(false)
			log.WithError(err).Error("Could not load pending cron jobs")
			time.Sleep(30 * time.Second)
			continue
		}
		c.recordPoll(true)

		for _, cronjob := range res {
			// Persist run to DB right away
			cronFields := log.Fields{"Cron (id)": cronjob.Cron, "type": cronjob.Type.String, "LastRun": cronjob.LastRun, "interval": cronjob.Interval, "id (not real id)": cronjob.ID.Int64}
			log.WithFields(cronFields).Trace("Running cronjob")

			claimed, err := c.claimJob(&cronjob)
			if err != nil {
				log.WithFields(cronFields).WithError(err).Error("Could not persist cron run")
				continue
			}
			if !claimed {
				log.WithFields(cronFields).Trace("Cron job changed or was claimed since polling")
				continue
			}

			// Log every failure with its job, since errgroup only returns the first.
			run := func(job func() error) {
				g.Go(func() error {
					c.recordJob(cronjob.Cron, false, nil)
					err := job()
					c.recordJob(cronjob.Cron, true, err)
					if err != nil {
						log.WithFields(cronFields).WithError(err).Error("Cron job failed")
					}
					return err
				})
			}
			// Run each job in a separate goroutine.
			switch cronjob.Type.String {
			case TUMOnlineStatusType:
				run(func() error { return c.tumonlineStatusCron() })
			case StudentClubType:
				// sequential, since both languages share the same images and would otherwise race on creating the file rows
				run(func() error {
					return errors.Join(c.studentClubCron(pb.Language_German), c.studentClubCron(pb.Language_English))
				})
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
