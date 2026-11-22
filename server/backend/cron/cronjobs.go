package cron

import (
	"context"
	"errors"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/getsentry/sentry-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const defaultPollInterval = 30 * time.Second

var (
	cronJobRunsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "campus_cron_job_runs_total",
		Help: "Total number of cron job executions by type and outcome.",
	}, []string{"type", "status"})
	cronJobDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "campus_cron_job_duration_seconds",
		Help:    "Duration of cron job executions by type.",
		Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60, 120, 300},
	}, []string{"type"})
)

type CronService struct {
	db           *gorm.DB
	client       *http.Client
	pollInterval time.Duration

	mu      sync.Mutex
	running map[int64]bool
	wg      sync.WaitGroup

	// executeOverride is used in tests to substitute cron job execution.
	executeOverride func(cronjob model.Crontab) error
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
		db:           db,
		client:       webfetch.Client,
		pollInterval: defaultPollInterval,
		running:      make(map[int64]bool),
	}
}

func (c *CronService) Run() error {
	return c.RunContext(context.Background())
}

// RunContext polls for due cron jobs until ctx is cancelled, then waits for
// any in-flight jobs to complete. Jobs run independently without blocking the
// scheduler loop, while a job that is still running will not overlap itself.
func (c *CronService) RunContext(ctx context.Context) error {
	interval := c.pollInterval
	if interval <= 0 {
		interval = defaultPollInterval
	}

	c.pollOnce(ctx)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("Cron: shutting down, waiting for running jobs to finish")
			c.wg.Wait()
			return nil
		case <-ticker.C:
			c.pollOnce(ctx)
		}
	}
}

func (c *CronService) pollOnce(ctx context.Context) {
	log.Trace("Cron: checking for pending")
	var res []model.Crontab

	if err := c.db.WithContext(ctx).Model(&model.Crontab{}).
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
		if ctx.Err() == nil {
			log.WithError(err).Error("Could not load pending cron jobs")
		}
		return
	}

	for _, cronjob := range res {
		if !c.tryMarkRunning(cronjob.Cron) {
			log.WithFields(log.Fields{"Cron (id)": cronjob.Cron, "type": cronjob.Type.String}).
				Warn("Skipping cron job because previous run is still in progress")
			continue
		}

		// Persist run to DB right away
		cronFields := log.Fields{
			"Cron (id)":        cronjob.Cron,
			"type":             cronjob.Type.String,
			"LastRun":          cronjob.LastRun,
			"interval":         cronjob.Interval,
			"id (not real id)": cronjob.ID.Int64,
		}
		log.WithFields(cronFields).Trace("Running cronjob")

		cronjob.LastRun = int32(time.Now().Unix())
		if err := c.db.WithContext(ctx).Save(&cronjob).Error; err != nil {
			c.unmarkRunning(cronjob.Cron)
			log.WithFields(cronFields).WithError(err).Error("Could not persist cron run")
			continue
		}

		jobCopy := cronjob
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			defer c.unmarkRunning(jobCopy.Cron)
			c.executeJob(jobCopy, cronFields)
		}()
	}
}

func (c *CronService) tryMarkRunning(cronID int64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.running == nil {
		c.running = make(map[int64]bool)
	}
	if c.running[cronID] {
		return false
	}
	c.running[cronID] = true
	return true
}

func (c *CronService) unmarkRunning(cronID int64) {
	c.mu.Lock()
	delete(c.running, cronID)
	c.mu.Unlock()
}

func (c *CronService) executeJob(cronjob model.Crontab, cronFields log.Fields) {
	jobType := cronjob.Type.String
	start := time.Now()
	defer func() {
		cronJobDuration.WithLabelValues(jobType).Observe(time.Since(start).Seconds())
		if r := recover(); r != nil {
			cronJobRunsTotal.WithLabelValues(jobType, "panic").Inc()
			log.WithFields(cronFields).WithFields(log.Fields{
				"panic": r,
				"stack": string(debug.Stack()),
			}).Error("Cron job panicked")
			sentry.CurrentHub().Recover(r)
		}
	}()

	var err error
	if c.executeOverride != nil {
		err = c.executeOverride(cronjob)
	} else {
		err = c.dispatchJob(cronjob)
	}
	if err != nil {
		cronJobRunsTotal.WithLabelValues(jobType, "error").Inc()
		log.WithFields(cronFields).WithError(err).Error("Cron job failed")
		return
	}
	cronJobRunsTotal.WithLabelValues(jobType, "ok").Inc()
}

func (c *CronService) dispatchJob(cronjob model.Crontab) error {
	switch cronjob.Type.String {
	case TUMOnlineStatusType:
		return c.tumonlineStatusCron()
	case StudentClubType:
		// sequential, since both languages share the same images and would otherwise race on creating the file rows
		return errors.Join(c.studentClubCron(pb.Language_German), c.studentClubCron(pb.Language_English))
	case NewsType:
		copyCronjob := cronjob
		return c.newsCron(&copyCronjob)
	case FileDownloadType:
		return c.fileDownloadCron()
	case DishNameDownload:
		return c.dishNameDownloadCron()
	case MovieType:
		return c.movieCron()
	case CanteenHeadcount:
		return c.canteenHeadCountCron()
	case FeedbackEmail:
		return c.feedbackEmailCron()
	default:
		return nil
	}
}
