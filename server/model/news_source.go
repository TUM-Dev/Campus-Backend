package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

// NewsSource struct is a row record of the newsSource table in the tca database
type NewsSource struct {
	Source  int32       `gorm:"primary_key;AUTO_INCREMENT;column:source;type:int;"`
	Title   string      `gorm:"column:title;type:text;size:16777215;"`
	URL     null.String `gorm:"column:url;type:text;size:16777215;"`
	FilesID int32       `gorm:"column:icon;not null"`
	Files   Files       `gorm:"foreignKey:FilesID;references:file;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Hook    null.String `gorm:"column:hook;type:char;size:12;"`
}

// TableName sets the insert table name for this struct type
func (n *NewsSource) TableName() string {
	return "newsSource"
}
