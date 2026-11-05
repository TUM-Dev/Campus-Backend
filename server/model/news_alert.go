package model

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

// NewsAlert struct is a row record of the news_alert table in the tca database
type NewsAlert struct {
	NewsAlert int64       `gorm:"primary_key;autoIncrement;column:news_alert;type:int;" json:"news_alert"`
	FileID    int64       `gorm:"column:file;not null;index:news_alert_files_file_fk"`
	File      File        `gorm:"foreignKey:FileID;references:file;constraint:-;"`
	Name      null.String `gorm:"column:name;type:varchar(100);" json:"name"`
	Link      null.String `gorm:"column:link;type:text;size:65535;" json:"link"`
	Created   time.Time   `gorm:"column:created;type:timestamp;not null;default:current_timestamp();" json:"created"`
	From      time.Time   `gorm:"column:from;type:datetime;not null;default:current_timestamp();" json:"from"`
	To        time.Time   `gorm:"column:to;type:datetime;not null;default:current_timestamp();" json:"to"`
}
