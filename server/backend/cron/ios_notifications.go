package cron

import (
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/device"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/devices_activity_reset"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/scheduled_update_log"
	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/scheduling"
)

// Starts the cron job for sending iOS notifications reuses
// the APNs client (apns.Service) stored in CronService
func (c *CronService) iosNotificationsCron() error {
	if !c.APNs.IsActive {
		return nil
	}

	devicesRepo := device.NewRepository(c.db)
	schedulerRepo := scheduled_update_log.NewRepository(c.db)

	service := scheduling.NewService(devicesRepo, schedulerRepo, c.APNs)

	return service.HandleScheduledCron()
}

// Resets the activity of all devices to 0 every day, week, month or year
func (c *CronService) iosActivityReset() error {
	service := devices_activity_reset.NewService(c.db)

	return service.HandleScheduledActivityReset()
}
