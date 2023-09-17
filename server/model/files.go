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

// Files struct is a row record of the files table in the tca database
type Files struct {
	File       int32       `gorm:"primary_key;AUTO_INCREMENT;column:file;type:int;" json:"file"`
	Name       string      `gorm:"column:name;type:text;size:16777215;" json:"name"`
	Path       string      `gorm:"column:path;type:text;size:16777215;" json:"path"`
	Downloads  int32       `gorm:"column:downloads;type:int;default:0;" json:"downloads"`
	URL        null.String `gorm:"column:url;default:null;" json:"url"`                         // URL of the files source (if any)
	Downloaded null.Bool   `gorm:"column:downloaded;type:boolean;default:1;" json:"downloaded"` // true when file is ready to be served, false when still being downloaded
}

// TableName sets the insert table name for this struct type
func (f *Files) TableName() string {
	return "files"
}
