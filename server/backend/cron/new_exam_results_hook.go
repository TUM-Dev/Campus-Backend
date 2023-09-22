package cron

import (
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/new_exam_results_hook/new_exam_results_scheduling"
)

func (c *CronService) newExamResultsHookCron() error {
	repo := new_exam_results_scheduling.NewRepository(c.db)
	devicesRepo := device.NewRepository(c.db)

	service := new_exam_results_scheduling.NewService(repo, devicesRepo, c.APNs)

	return service.HandleScheduledCron()
}
