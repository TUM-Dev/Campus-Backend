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

type NewsSource struct {
	Source int64       `gorm:"primary_key;autoIncrement;column:source;type:int;"`
	Title  string      `gorm:"column:title;type:text;not null;"`
	URL    null.String `gorm:"column:url;type:text;size:16777215;"`
	FileID int64       `gorm:"column:icon;not null;type:int;index:icon;"`
	File   File        `gorm:"foreignKey:FileID;references:file;constraint:-;"`
	Hook   null.String `gorm:"column:hook;type:enum('newspread','impulsivHook');"`
}
