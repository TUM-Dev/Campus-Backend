package model

import (
	"fmt"
	"time"
)

// IOSSchedulingPriority stores some default priorities for the scheduling of
// grade updates.
type IOSSchedulingPriority struct {
	ID       int64 `gorm:"primary_key;auto_increment;not_null" json:"id"`
	FromDay  int   `gorm:"not null" json:"from_day"`
	ToDay    int   `gorm:"not null" json:"to_day"`
	FromHour int   `gorm:"not null" json:"from_hour"`
	ToHour   int   `gorm:"not null" json:"to_hour"`
	Priority int   `gorm:"not null" json:"priority"`
}

// IsCurrentlyInRange returns true if the current time is in the range of the
// scheduling priority.
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

// IsMorePreciseThan compares two Priorities and returns true if the current
// priority is more precise than the other one.
//
// Example:
// A priority with FromDay=1, ToDay=365, FromHour=6, ToHour=8 is more precise
// than a priority with FromDay=1, ToDay=365, FromHour=0, ToHour=23.
// In case the current hour is between 6 and 8.
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
