package cron

import (
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CronService struct {
	db   *gorm.DB
	cron *cron.Cron // robfig/cron
}

func New(db *gorm.DB) *CronService {
	return &CronService{
		db:   db,
		cron: cron.New(),
	}
}

func (c *CronService) Run() error {

	return nil
}
