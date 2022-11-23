package cron

import (
	"github.com/TUM-Dev/Campus-Backend/backend/ios_notifications/ios_scheduling"
)

func (c *CronService) iosNotificationsCron() error {
	repo := ios_scheduling.NewRepository(c.db)

	service := ios_scheduling.NewService(repo)

	return service.HandleScheduledCron()
}
