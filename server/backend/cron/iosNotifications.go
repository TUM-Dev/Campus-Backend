package cron

import (
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_scheduled_update_log"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_scheduling"
)

// Starts the cron job for sending iOS notifications reuses
// the apns client (ios_apns.Service) stored in CronService
func (c *CronService) iosNotificationsCron() error {
	repo := ios_scheduling.NewRepository(c.db)
	devicesRepo := ios_device.NewRepository(c.db)
	schedulerRepo := ios_scheduled_update_log.NewRepository(c.db)

	service := ios_scheduling.NewService(repo, devicesRepo, schedulerRepo, c.APNs)

	return service.HandleScheduledCron()
}
