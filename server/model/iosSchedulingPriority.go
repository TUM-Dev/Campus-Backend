package model

import (
	"fmt"
	"time"
)

type IOSSchedulingPriority struct {
	ID       int `gorm:"primary_key;auto_increment;not_null" json:"id"`
	FromDay  int `gorm:"not null" json:"from_day"`
	ToDay    int `gorm:"not null" json:"to_day"`
	FromHour int `gorm:"not null" json:"from_hour"`
	ToHour   int `gorm:"not null" json:"to_hour"`
	Priority int `gorm:"not null" json:"priority"`
}

func (p *IOSSchedulingPriority) IsCurrentlyInRange() bool {
	now := time.Now()
	yearDay := now.YearDay()

	if p.FromDay <= yearDay && p.ToDay >= yearDay {
		hour := now.Hour()

		if p.FromHour <= hour && p.ToHour >= hour {
			return true
		}
	}

	return false
}

func (p *IOSSchedulingPriority) IsMorePreciseThan(other *IOSSchedulingPriority) bool {
	if p.FromDay == other.FromDay && p.ToDay == other.ToDay {
		if p.FromHour == other.FromHour && p.ToHour == other.ToHour {
			return p.Priority > other.Priority
		}

		return p.FromHour > other.FromHour && p.ToHour < other.ToHour
	}

	return p.FromDay > other.FromDay && p.ToDay < other.ToDay
}

func (p *IOSSchedulingPriority) String() string {
	return fmt.Sprintf("Day: %d-%d, Hour: %d-%d, Priority: %d", p.FromDay, p.ToDay, p.FromHour, p.ToHour, p.Priority)
}

func DefaultIOSSchedulingPriority() *IOSSchedulingPriority {
	return &IOSSchedulingPriority{
		FromDay:  1,
		ToDay:    365,
		FromHour: 0,
		ToHour:   23,
		Priority: 5,
	}
}
