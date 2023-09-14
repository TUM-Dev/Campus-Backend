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
	News        int32       `gorm:"primary_key;AUTO_INCREMENT;column:news;type:int;"`
	Date        time.Time   `gorm:"column:date;type:datetime;"`
	Created     time.Time   `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;"`
	Title       string      `gorm:"column:title;type:text;size:255;"`
	Description string      `gorm:"column:description;type:text;size:65535;"`
	Src         int32       `gorm:"column:src;type:int;"`
	Link        string      `gorm:"column:link;type:varchar(190);"`
	Image       null.String `gorm:"column:image;type:text;size:65535;"`
	FilesID     null.Int    `gorm:"column:file;type:int;"`
	Files       Files       `gorm:"foreignKey:FilesID;references:file;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName sets the insert table name for this struct type
func (n *News) TableName() string {
	return "news"
}
