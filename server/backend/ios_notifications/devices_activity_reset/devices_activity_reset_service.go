package devices_activity_reset

import (
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/backend/ios_notifications/device"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Service struct {
	Repository *Repository
}

func (service *Service) HandleScheduledActivityReset() error {
	daily, err := service.Repository.GetDevicesActivityResetDaily()

	if err != nil {
		service.Repository.CreateInitialRecords()

		return nil
	}

	weekly, err := service.Repository.GetDevicesActivityResetWeekly()

	if err != nil {
		return err
	}

	monthly, err := service.Repository.GetDevicesActivityResetMonthly()

	if err != nil {
		return err
	}

	yearly, err := service.Repository.GetDevicesActivityResetYearly()

	if err != nil {
		return err
	}

	now := time.Now()

	devicesRepo := device.NewRepository(service.Repository.DB)

	if now.Sub(daily.LastReset).Hours() > 24 {
		if err := service.Repository.ResettedDevicesDaily(); err != nil {
			log.WithError(err).Error("while resetting devices daily")
		}

		if err := devicesRepo.ResetDevicesDailyActivity(); err != nil {
			log.WithError(err).Error("while resetting devices daily activity")
		}
	}

	if now.Sub(weekly.LastReset).Hours() > 168 {
		if err := service.Repository.ResettedDevicesWeekly(); err != nil {
			log.WithError(err).Error("while resetting devices weekly")
		}

		if err := devicesRepo.ResetDevicesWeeklyActivity(); err != nil {
			log.WithError(err).Error("while resetting devices weekly activity")
		}
	}

	if now.Sub(monthly.LastReset).Hours() > 730 {
		if err := service.Repository.ResettedDevicesMonthly(); err != nil {
			log.WithError(err).Error("while resetting devices monthly")
		}

		if err := devicesRepo.ResetDevicesMonthlyActivity(); err != nil {
			log.WithError(err).Error("while resetting devices monthly activity")
		}
	}

	if now.Sub(yearly.LastReset).Hours() > 8760 {
		if err := service.Repository.ResettedDevicesYearly(); err != nil {
			log.WithError(err).Error("while resetting devices yearly")
		}

		if err := devicesRepo.ResetDevicesYearlyActivity(); err != nil {
			log.WithError(err).Error("while resetting devices yearly activity")
		}
	}

	return nil
}

func NewService(db *gorm.DB) *Service {
	return &Service{Repository: NewRepository(db)}
}
