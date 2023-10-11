package cron

import (
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/new_exam_results_hook/scheduling"
)

func (c *CronService) newExamResultsHookCron() error {
	repo := scheduling.NewRepository(c.db)
	devicesRepo := device.NewRepository(c.db)

	service := scheduling.NewService(repo, devicesRepo, c.APNs)

	return service.HandleScheduledCron()
}
