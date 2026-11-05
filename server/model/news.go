package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

// News struct is a row record of the news table in the tca database
type News struct {
	News         int64       `gorm:"primary_key;autoIncrement;column:news;type:int;"`
	Date         time.Time   `gorm:"column:date;type:datetime;not null;"`
	Created      time.Time   `gorm:"column:created;type:timestamp;not null;default:current_timestamp();"`
	Title        string      `gorm:"column:title;type:tinytext;not null;"`
	Description  string      `gorm:"column:description;type:text;not null;"`
	NewsSourceID int64       `gorm:"column:src;type:int;not null;index:src;"`
	NewsSource   NewsSource  `gorm:"foreignKey:NewsSourceID;references:source;constraint:-;"`
	Link         string      `gorm:"column:link;type:varchar(190);not null;uniqueIndex:link;"`
	Image        null.String `gorm:"column:image;type:text;size:65535;"`
	FileID       null.Int    `gorm:"column:file;type:int;index:file;"`
	File         *File       `gorm:"foreignKey:FileID;references:file;constraint:-;"`
}
