package model

import (
	"github.com/guregu/null"
	"time"
)

// Movie stores all movies
type Movie struct {
	Id          int64       `gorm:"primary_key;AUTO_INCREMENT;column:kino;type:int;not null;"`
	Date        time.Time   `gorm:"column:date;type:datetime;not null;"`
	Created     time.Time   `gorm:"column:created;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Title       string      `gorm:"column:title;type:text;not null;"`
	Year        null.String `gorm:"column:year;type:varchar(4)"`
	Runtime     null.String `gorm:"column:runtime;type:varchar(40)"`
	Genre       null.String `gorm:"column:genre;type:varchar(100)"`
	Director    null.String `gorm:"column:director;type:text"`
	Actors      null.String `gorm:"column:actors;type:text"`
	ImdbRating  null.String `gorm:"column:rating;type:varchar(4)"`
	Description string      `gorm:"column:description;type:text;not null;"`
	Trailer     null.String `gorm:"column:trailer"`
	FileID      int64       `gorm:"column:cover"`
	File        File        `gorm:"foreignKey:FileID;references:file"`
	Link        string      `gorm:"column:link;type:varchar(190);not null;unique;"`
	Location    null.String `gorm:"column:location;default:null"`
}
