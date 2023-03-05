package ios_device

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"gorm.io/gorm"
)

const (
	DevicesThatCouldUpdateLectureCount = 10
)

type Repository struct {
	DB *gorm.DB
}

func (repository *Repository) RegisterDevice(device *model.IOSDevice) (bool, error) {
	exists, err := repository.CheckIfDeviceExists(device.DeviceID)

	if err != nil {
		return false, err
	}

	err = repository.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.
			Where("device_id = ?", device.DeviceID).
			FirstOrCreate(&device).Error

		if err != nil {
			return err
		}

		newDevice := model.IOSDevice{
			DeviceID:          device.DeviceID,
			PublicKey:         device.PublicKey,
			ActivityToday:     device.ActivityToday + 1,
			ActivityThisWeek:  device.ActivityThisWeek + 1,
			ActivityThisMonth: device.ActivityThisMonth + 1,
			ActivityThisYear:  device.ActivityThisYear + 1,
		}

		return tx.Save(&newDevice).Error
	})

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (repository *Repository) CheckIfDeviceExists(deviceId string) (bool, error) {
	var devices []model.IOSDevice
	if err := repository.DB.Limit(1).Find(&devices, "device_id = ?", deviceId).Error; err != nil {
		return false, err
	}

	return len(devices) > 0, nil
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

func (repository *Repository) GetDevice(id string) (*model.IOSDevice, error) {
	var device *model.IOSDevice
	if err := repository.DB.First(&device, "device_id = ?", id).Error; err != nil {
		return nil, err
	}

	return device, nil
}

// GetDevicesThatShouldUpdateGrades returns a list of devices that should be updated
// A device that needs to be updated is either a new device or a device that has not
// been updated in the last model.IOSMinimumUpdateInterval minutes
func (repository *Repository) GetDevicesThatShouldUpdateGrades() ([]model.IOSDeviceLastUpdated, error) {
	var devices []model.IOSDeviceLastUpdated

	tx := repository.DB.Raw(
		buildDevicesThatShouldUpdateGradesQuery(),
		model.IOSUpdateTypeGrades,
		model.IOSMinimumUpdateInterval,
	).Scan(&devices)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return devices, nil
}

func buildDevicesThatShouldUpdateGradesQuery() string {
	return `
		select d.device_id, ul.created_at as last_updated, d.public_key
		from ios_devices d
				 left join ios_scheduled_update_logs ul on d.device_id = ul.device_id
		where ul.created_at is null
		   or (ul.type = ?
			and ul.created_at < date_sub(now(), interval ? minute))
		group by d.device_id, ul.created_at
		order by ul.created_at;
	`
}

func (repository *Repository) GetDevicesThatCouldUpdateLecture(lecture *model.IOSLecture) ([]model.IOSDevice, error) {
	var devices []model.IOSDevice

	tx := repository.DB.Raw(
		buildGetDevicesThatCouldUpdateLectureQuery(),
		lecture.Id,
		DevicesThatCouldUpdateLectureCount,
	).Scan(&devices)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return devices, nil
}

func buildGetDevicesThatCouldUpdateLectureQuery() string {
	return `
	with ios_devices_with_lecture as (select d.*
                                  from ios_devices d,
                                       ios_device_lectures dl
                                  where d.device_id = dl.device_id
                                    and dl.lecture_id = ?),
     ios_devices_response_time
         as (select avg(timestampdiff(SECOND, drl.created_at, drl.handled_at)) as avg_response_time, drl.device_id
             from ios_device_request_logs drl,
                  ios_devices_with_lecture dl
             where dl.device_id = drl.device_id
               and drl.handled_at is not null),
     ios_devices_lecture_count as (select count(device_id) as lecture_count, device_id
                                   from ios_device_lectures
                                   group by device_id)
	select d.*
	from ios_devices d
			 left join ios_devices_response_time drt on d.device_id = drt.device_id
				left join ios_devices_lecture_count dlc on d.device_id = dlc.device_id,
		 ios_device_lectures dl
	where d.device_id = dl.device_id
	  and not exists(select d.device_id
					 from ios_device_request_logs drl
					 where drl.device_id = d.device_id
					   and drl.created_at
						 > subdate(now()
							   , interval 30 minute)
					 limit 1)
	group by d.device_id
	order by dlc.lecture_count desc, drt.avg_response_time asc
	limit ?;
`
}

func (repository *Repository) ResetDevicesDailyActivity() error {
	return repository.DB.Model(model.IOSDevice{}).Where("activity_today != ?", 0).Update("activity_today", 0).Error
}

func (repository *Repository) ResetDevicesWeeklyActivity() error {
	return repository.DB.Model(model.IOSDevice{}).Where("activity_this_week != ?", 0).Update("activity_this_week", 0).Error
}

func (repository *Repository) ResetDevicesMonthlyActivity() error {
	return repository.DB.Model(model.IOSDevice{}).Where("activity_this_month != ?", 0).Update("activity_this_month", 0).Error
}

func (repository *Repository) ResetDevicesYearlyActivity() error {
	return repository.DB.Model(model.IOSDevice{}).Where("activity_this_year != ?", 0).Update("activity_this_year", 0).Error
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
