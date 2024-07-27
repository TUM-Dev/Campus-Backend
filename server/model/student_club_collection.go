package model

import (
	"time"

	"gorm.io/gorm"
)

// StudentClubCollection stores what collection a club belongs to
type StudentClubCollection struct {
	ID          string `gorm:"primaryKey;type:varchar(100)"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
