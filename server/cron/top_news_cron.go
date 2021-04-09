package cron

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	"log"
)

func (c ServiceCron) fetchNews() {
	// do some networking stuff, load news
	var newsEntry = model.TopNews{
		Name:    "newsTest",
		Link:    "https://tum.de",
		Created: nil,
		From:    nil,
		To:      nil,
	}
	// add them to the database
	id := c.DB.Create(&newsEntry)
	log.Printf("created news with id %v\n", *id)
}
