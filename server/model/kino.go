package model

import (
	"time"

	"github.com/guregu/null"
)

// Kino stores all movies
type Kino struct {
	Id          int64       `gorm:"primary_key;AUTO_INCREMENT;column:kino;type:int;not null;"`
	Date        time.Time   `gorm:"column:date;type:datetime;not null;"`
	Created     time.Time   `gorm:"column:created;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Title       string      `gorm:"column:title;type:text;not null;"`
	Year        string      `gorm:"column:year;type:varchar(4);not null;"`
	Runtime     string      `gorm:"column:runtime;type:varchar(40);not null;"`
	Genre       string      `gorm:"column:genre;type:varchar(100);not null;"`
	Director    string      `gorm:"column:director;type:text;not null;"`
	Actors      string      `gorm:"column:actors;type:text;not null;"`
	ImdbRating  string      `gorm:"column:rating;type:varchar(4);not null;"`
	Description string      `gorm:"column:description;type:text;not null;"`
	Trailer     null.String `gorm:"column:trailer"`
	FileID      int64       `gorm:"column:cover"`
	File        File        `gorm:"foreignKey:FileID;references:file"`
	Link        string      `gorm:"column:link;type:varchar(190);not null;unique;"`
}

// TableName sets the insert table name for this struct type
func (n *Kino) TableName() string {
	return "kino"
}
