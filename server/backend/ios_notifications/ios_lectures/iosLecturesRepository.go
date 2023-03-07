package ios_lectures

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	DB *gorm.DB
}

// SaveLectureAsIOSLecture saves a lecture as an IOSLecture => filtering out unnecessary fields
func (repo *Repository) SaveLectureAsIOSLecture(lecture *model.Lecture) error {
	iosLecture, err := lecture.ToIOSLecture()
	if err != nil {
		return err
	}

	if err := repo.DB.Create(&iosLecture).Error; err != nil {
		return err
	}

	return nil
}

// SaveLecturesAsIOSLectures saves a list of lectures as IOSLectures => batching the insertions
func (repo *Repository) SaveLecturesAsIOSLectures(lectures []model.Lecture) error {
	var iosLectures []model.IOSLecture

	for _, lecture := range lectures {
		iosLecture, err := lecture.ToIOSLecture()
		if err != nil {
			continue
		}

		iosLectures = append(iosLectures, *iosLecture)
	}

	if len(iosLectures) == 0 {
		log.Info("No lectures to save")
		return nil
	}

	tx := repo.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&iosLectures)
	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func (repo *Repository) SaveLecturesOfDevice(lectures []model.Lecture, deviceId string) error {
	var deviceLectures []model.IOSDeviceLecture

	for _, lecture := range lectures {
		deviceLecture := model.IOSDeviceLecture{
			DeviceId:  deviceId,
			LectureId: lecture.Id,
		}

		deviceLectures = append(deviceLectures, deviceLecture)
	}

	tx := repo.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_id"}, {Name: "lecture_id"}},
		DoNothing: true,
	}).Create(&deviceLectures)

	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func (repo *Repository) GetLectures() (*[]model.IOSLecture, error) {
	var lectures []model.IOSLecture

	if err := repo.DB.Find(&lectures).Error; err != nil {
		return nil, err
	}

	return &lectures, nil
}

func (repo *Repository) GetDeviceLectures() (*[]model.IOSDeviceLecture, error) {
	var deviceLectures []model.IOSDeviceLecture

	if err := repo.DB.Model(&model.IOSDeviceLecture{}).Find(&deviceLectures).Error; err != nil {
		return nil, err
	}

	return &deviceLectures, nil
}

// SetLecturesLastUpdatedBy sets the last_request_id of all lectures that are in the given grades
// IMPORTANT: the title needs to be filtered with `like` because the title of the lecture in the grades
// is not exactly the same as the title of the lecture in the ios_lectures table.
func (repo *Repository) SetLecturesLastUpdatedBy(requestId string, grades *[]model.Grade) {
	for _, grade := range *grades {
		tx := repo.DB.Model(&model.IOSLecture{}).
			Where("title like ?", "%"+grade.LectureTitle+"%").
			Update("last_request_id", requestId)

		if err := tx.Error; err != nil {
			log.WithError(err).Error("Failed to update last_request_id for lecture")
		}
	}
}

// FindOtherDevicesThatAttendLecture finds all devices that attend a lecture with the given title which
// then can be notified about e.g. a new grade.
// IMPORTANT: the title needs to be filtered with `like` => see SetLecturesLastUpdatedBy
func (repo *Repository) FindOtherDevicesThatAttendLecture(lectureTitle string) (*[]model.IOSDevice, error) {
	var devices []model.IOSDevice

	tx := repo.DB.Raw(
		buildFindOtherDevicesThatAttendLectureQuery(),
		"%"+lectureTitle+"%",
	).Scan(&devices)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &devices, nil
}

func buildFindOtherDevicesThatAttendLectureQuery() string {
	return `
	select d.*
	from ios_lectures l,
		 ios_device_lectures dl,
		 ios_devices d
	where l.id = dl.lecture_id
	  and d.device_id = dl.device_id
	  and l.title like ?;
	`
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
