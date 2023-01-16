package ios_device

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	DB *gorm.DB
}

func (repository *Repository) RegisterDevice(device *model.IOSDevice) error {
	if err := repository.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"public_key"}),
	}).Create(device).Error; err != nil {
		return err
	}

	return nil
}

func (repository *Repository) RemoveDevice(deviceId string) error {
	if err := repository.DB.Delete(&model.IOSDevice{DeviceID: deviceId}).Error; err != nil {
		return err
	}

	return nil
}

func (repository *Repository) GetDevices() ([]model.IOSDevice, error) {
	var devices []model.IOSDevice
	if err := repository.DB.Find(&devices).Error; err != nil {
		return []model.IOSDevice{}, err
	}

	return devices, nil
}

func (repository *Repository) GetDevicesThatShouldUpdateGrades() ([]model.IOSDeviceLastUpdated, error) {
	var devices []model.IOSDeviceLastUpdated

	tx := repository.DB.Raw(repository.buildDevicesThatShouldUpdateGradesQuery(), model.IOSMinimumUpdateInterval).Scan(&devices)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return devices, nil
}

func (repository *Repository) buildDevicesThatShouldUpdateGradesQuery() string {
	return `
		select d.device_id, ul.created_at as last_updated, d.public_key
		from ios_devices d
				 left join ios_scheduled_update_logs ul on d.device_id = ul.device_id
		where ul.created_at is null
		   or (ul.type = 'grades'
			and ul.created_at < date_sub(now(), interval ? minute))
		group by d.device_id, ul.created_at
		order by ul.created_at;
	`
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
