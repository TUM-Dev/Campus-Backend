package model

type DeviceExam struct {
	DeviceId string    `gorm:"primaryKey"`
	Device   IOSDevice `gorm:"constraint:OnDelete:CASCADE"`
	ExamId   string    `gorm:"primaryKey"`
	Exam     DbExam    `gorm:"constraint:OnDelete:CASCADE"`
}
