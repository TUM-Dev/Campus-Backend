package ios_devices_activity_reset

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Repository struct {
	DB *gorm.DB
}

func (repo *Repository) GetDevicesActivityResets() ([]model.IOSDevicesActivityReset, error) {
	var resets []model.IOSDevicesActivityReset

	err := repo.DB.Find(&resets).Error

	if err != nil {
		return nil, err
	}

	return resets, nil
}

func (repo *Repository) GetDevicesActivityResetDaily() (*model.IOSDevicesActivityReset, error) {
	return repo.GetDevicesActivityReset(model.IOSDevicesActivityResetTypeDay)
}

func (repo *Repository) GetDevicesActivityResetWeekly() (*model.IOSDevicesActivityReset, error) {
	return repo.GetDevicesActivityReset(model.IOSDevicesActivityResetTypeWeek)
}

func (repo *Repository) GetDevicesActivityResetMonthly() (*model.IOSDevicesActivityReset, error) {
	return repo.GetDevicesActivityReset(model.IOSDevicesActivityResetTypeMonth)
}

func (repo *Repository) GetDevicesActivityResetYearly() (*model.IOSDevicesActivityReset, error) {
	return repo.GetDevicesActivityReset(model.IOSDevicesActivityResetTypeYear)
}

func (repo *Repository) GetDevicesActivityReset(resetType string) (*model.IOSDevicesActivityReset, error) {
	var reset model.IOSDevicesActivityReset
	if err := repo.DB.First(&reset, "type = ?", resetType).Error; err != nil {
		return nil, err
	}

	return &reset, nil
}

func (repo *Repository) ResettedDevicesDaily() error {
	return repo.ResetActivityNow(model.IOSDevicesActivityResetTypeDay)
}

func (repo *Repository) ResettedDevicesWeekly() error {
	return repo.ResetActivityNow(model.IOSDevicesActivityResetTypeWeek)
}

func (repo *Repository) ResettedDevicesMonthly() error {
	return repo.ResetActivityNow(model.IOSDevicesActivityResetTypeMonth)
}

func (repo *Repository) ResettedDevicesYearly() error {
	return repo.ResetActivityNow(model.IOSDevicesActivityResetTypeYear)
}

func (repo *Repository) ResetActivityNow(resetType string) error {
	reset := model.IOSDevicesActivityReset{
		Type:      resetType,
		LastReset: time.Now(),
	}

	return repo.ResetActivity(&reset)
}

func (repo *Repository) ResetActivity(reset *model.IOSDevicesActivityReset) error {
	res := repo.DB.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "type"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"last_reset",
			}),
		},
	).Create(&reset)

	return res.Error
}

func (repo *Repository) CreateInitialRecords() {
	now := time.Now()

	types := []string{
		model.IOSDevicesActivityResetTypeDay,
		model.IOSDevicesActivityResetTypeWeek,
		model.IOSDevicesActivityResetTypeMonth,
		model.IOSDevicesActivityResetTypeYear,
	}

	// iterate over types
	for _, resetType := range types {
		reset := model.IOSDevicesActivityReset{
			Type:      resetType,
			LastReset: now,
		}

		if err := repo.DB.Create(&reset).Error; err != nil {
			log.WithError(err).WithField("resetType", resetType).Error("Failed to create IOSDevicesActivityReset")
			continue
		}
	}
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
