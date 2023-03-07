package cron

import (
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_devices_activity_reset"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/ios_scheduling"
)

// Starts the cron job for sending iOS notifications reuses
// the APNs client (ios_apns.Service) stored in CronService
func (c *CronService) iosNotificationsCron() error {
	if !c.APNs.IsActive {
		return nil
	}

	repo := ios_scheduling.NewRepository(c.db)
	devicesRepo := ios_device.NewRepository(c.db)

	service := ios_scheduling.NewService(repo, devicesRepo, c.APNs)

	return service.HandleScheduledCron()
}

// Resets the activity of all devices to 0 every day, week, month or year
func (c *CronService) iosActivityReset() error {
	service := ios_devices_activity_reset.NewService(c.db)

	return service.HandleScheduledActivityReset()
}
