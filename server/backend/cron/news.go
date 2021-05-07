package cron

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/getsentry/sentry-go"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
)

// newsCron fetches news and saves them to the database
func (c *CronService) newsCron() error {
	var newsSources []model.NewsSource
	err := c.db.Find(&newsSources).Error
	if err != nil {
		log.Printf("error getting news sources from database: %v", err)
		sentry.CaptureException(err)
	}
	for _, source := range newsSources {
		// skip sources with null url.
		if source.URL.Valid {
			log.Println("processing source %v", source.URL)
			feed, err := c.gf.ParseURL(source.URL.String)
			if err != nil {
				log.Printf("error parsing rss: %v", err)
				sentry.CaptureException(err)
			} else {
				c.parseNewsFeed(feed)
			}
		}
	}
	return nil
}

func (c *CronService) parseNewsFeed(feed *gofeed.Feed) {
	log.Println(feed.Title)
}
