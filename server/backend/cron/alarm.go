package cron

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (c *CronService) alarmCron() error {
	return notifyAllUsers(c.db, sendGCMNotification)
}

func sendGCMNotification(*[]model.Devices) error {
	log.Trace("sendGCMNotification")
	return nil
}

func notifyAllUsers(db *gorm.DB, callbacks ...func(*[]model.Devices) error) error {
	var pendingNotifications []model.Notification
	if err := db.
		Joins("notification_confirmation").
		Where("notification_confirmation.sent = ?", false).
		Group("notification").
		Find(&pendingNotifications).Error; err != nil {
		log.WithError(err).Error("failed to get pending notifications from the db")
		return err
	}
	for _, pending := range pendingNotifications {
		// Get a few targets
		var targets []model.Devices
		if err := db.
			Distinct("device").
			Where("notification = ?", pending.Notification).
			Where("gcmStatus IS NULL").
			Where("gcmToken IS NOT NULL").
			Joins("notification_confirmation").
			Where("notification_confirmation.sent = ?", false).
			Limit(998).
			Find(&targets).
			Error; err != nil {
			log.WithError(err).Error("failed to get devices which should receive the notification")
			continue
		}
		if len(targets) > 0 {
			for _, callback := range callbacks {
				if err := callback(&targets); err != nil {
					log.WithError(err).Error("callback failed")
					continue
				}
			}
			// mark as sent for these targets
			if err := db.
				Model(&model.NotificationConfirmation{}).
				Where("notification = ?", pending.Notification).
				Where("device IN ?", targets).
				Update("sent", true).
				Error; err != nil {
				log.WithError(err).Error("failed to update notification_confirmation")
				continue
			}
		}
	}
	return nil
}
