package model

type IOSDeviceLecture struct {
	DeviceId  string     `gorm:"primaryKey"`
	Device    IOSDevice  `gorm:"constraint:OnDelete:CASCADE"`
	LectureId string     `gorm:"primaryKey"`
	Lecture   IOSLecture `gorm:"constraint:OnDelete:CASCADE"`
}
