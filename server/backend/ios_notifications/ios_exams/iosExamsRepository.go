package ios_exams

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	DB *gorm.DB
}

func (repo *Repository) SaveExams(exams []model.Exam) error {
	var dbExams []model.DbExam

	for _, exam := range exams {
		dbExams = append(dbExams, exam.ToDbExam())
	}

	if len(dbExams) == 0 {
		log.Info("No lectures to save")
		return nil
	}

	tx := repo.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&dbExams)
	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func (repo *Repository) SaveExamsOfDevice(exams []model.Exam, deviceId string) error {
	var deviceExams []model.DeviceExam

	for _, exam := range exams {
		deviceExam := model.DeviceExam{
			DeviceId: deviceId,
			ExamId:   exam.ExamID,
		}

		deviceExams = append(deviceExams, deviceExam)
	}

	tx := repo.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_id"}, {Name: "exam_id"}},
		DoNothing: true,
	}).Create(&deviceExams)

	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func (repo *Repository) GetExams() (*[]model.DbExam, error) {
	var exams []model.DbExam

	if err := repo.DB.Find(&exams).Error; err != nil {
		return nil, err
	}

	return &exams, nil
}

func (repo *Repository) GetDeviceExams() (*[]model.DeviceExam, error) {
	var deviceExams []model.DeviceExam

	if err := repo.DB.Model(&model.DeviceExam{}).Find(&deviceExams).Error; err != nil {
		return nil, err
	}

	return &deviceExams, nil
}

func (repo *Repository) GetDevicesThatHaveExams(examIds *[]string) (*[]model.DeviceExam, error) {
	var deviceExams []model.DeviceExam

	if err := repo.DB.Model(&model.DeviceExam{}).Where("exam_id IN ?", *examIds).Find(&deviceExams).Error; err != nil {
		return nil, err
	}

	return &deviceExams, nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		DB: db,
	}
}
