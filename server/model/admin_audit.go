package model

import "time"

type AdminAudit struct {
	ID        uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime(6);not null"`
	Subject   string    `gorm:"column:subject;type:varchar(255);not null"`
	Action    string    `gorm:"column:action;type:varchar(32);not null"`
	SourceID  int64     `gorm:"column:source_id;type:bigint;not null"`
	Details   string    `gorm:"column:details;type:mediumtext;not null"`
}
