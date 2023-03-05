package ios_lectures

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	DB *gorm.DB
}

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

func (repo *Repository) SaveLecturesAsIOSLectures(lectures []model.Lecture) error {
	var iosLectures []model.IOSLecture

	for _, lecture := range lectures {
		iosLecture, err := lecture.ToIOSLecture()
		if err != nil {
			continue
		}

		iosLectures = append(iosLectures, *iosLecture)
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

func (repo *Repository) GetLecturesToUpdate() ([]model.IOSLecture, error) {
	var lectures []model.IOSLecture

	tx := repo.DB.Raw(
		buildGetLecturesToUpdateQuery(),
	).Scan(&lectures)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return lectures, nil
}

func buildGetLecturesToUpdateQuery() string {
	return `
		select lec.id, lec.year, lec.semester, lec.last_request_id
		from ios_lectures lec
				 left join ios_device_request_logs log on lec.last_request_id = log.request_id
		where lec.last_request_id is null
		   or (
			log.handled_at < subdate(now(), interval 1 minute)
			)
		order by log.handled_at asc;
	`
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
