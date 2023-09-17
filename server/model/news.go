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
	News        int32       `gorm:"primary_key;AUTO_INCREMENT;column:news;type:int;" json:"news"`
	Date        time.Time   `gorm:"column:date;type:datetime;" json:"date"`
	Created     time.Time   `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	Title       string      `gorm:"column:title;type:text;size:255;" json:"title"`
	Description string      `gorm:"column:description;type:text;size:65535;" json:"description"`
	Src         int32       `gorm:"column:src;type:int;" json:"src"`
	Link        string      `gorm:"column:link;type:varchar(190);" json:"link"`
	Image       null.String `gorm:"column:image;type:text;size:65535;" json:"image"`
	File        null.Int    `gorm:"column:file;type:int;" json:"file"`
}

// TableName sets the insert table name for this struct type
func (n *News) TableName() string {
	return "news"
}
