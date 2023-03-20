package cron

import (
	"github.com/TUM-Dev/Campus-Backend/server/backend/lecture_crawler"
	log "github.com/sirupsen/logrus"
)

func (c *CronService) lectureCrawlerCron() error {
	log.Infof("Running lecture crawler cron job")

	crawler := lecture_crawler.New(c.db)
	crawler.Crawl()

	return nil
}
