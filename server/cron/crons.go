package cron

import (
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type ServiceCron struct {
	DB *gorm.DB
}

func (c ServiceCron) Init() {
	cabeCron := cron.New()
	// fetch news once per hour
	_, _ = cabeCron.AddFunc("0 * * * *", func() { c.fetchNews() })
	cabeCron.Start()
}
