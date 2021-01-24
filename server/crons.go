package server

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
)

type CronService struct {
	DB *gorm.DB
}

func (c CronService) Init() {
	cabeCron := cron.New()
	// fetch news once per hour
	_, _ = cabeCron.AddFunc("0 * * * *", func() { fetchNews(c.DB) })
	cabeCron.Start()
}

func fetchNews(db *gorm.DB) {
	// do some networking stuff, load news
	var newsEntry = model.TopNews{
		Name:    "newsTest",
		Link:    "https://tum.de",
		Created: nil,
		From:    nil,
		To:      nil,
	}
	// add them to the database
	id := db.Create(&newsEntry)
	log.Printf("created news with id %v\n", *id)
}
