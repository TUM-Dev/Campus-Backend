package model

type DeviceLecture struct {
	DeviceId  string    `gorm:"primaryKey"`
	Device    IOSDevice `gorm:"constraint:OnDelete:CASCADE"`
	LectureId string    `gorm:"primaryKey"`
	Lecture   Lecture   `gorm:"constraint:OnDelete:CASCADE"`
}
