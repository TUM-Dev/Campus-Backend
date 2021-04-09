package cron

import (
	"github.com/getsentry/sentry-go"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type ServiceCron struct {
	DB         *gorm.DB
	CabeSentry sentry.Client
}

func (c ServiceCron) Init() {
	cabeCron := cron.New()
	// fetch news once per hour
	_, _ = cabeCron.AddFunc("0 * * * *", func() { c.fetchNews() })
	_, _ = cabeCron.AddFunc("* * * * *", func() { c.fetchRooms() })
	cabeCron.Start()
}
