package model

import (
	"database/sql"
	"fmt"
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

// File struct is a row record of the files table in the tca database
type File struct {
	File       int64       `gorm:"primary_key;AUTO_INCREMENT;column:file;type:int;" json:"file"`
	Name       string      `gorm:"column:name;type:text;size:16777215;not null" json:"name"`
	Path       string      `gorm:"column:path;type:text;size:16777215;not null" json:"path"`
	Downloads  int32       `gorm:"column:downloads;type:int;default:0;not null" json:"downloads"`
	URL        null.String `gorm:"column:url;default:null;" json:"url"`                         // URL of the files source (if any)
	Downloaded null.Bool   `gorm:"column:downloaded;type:boolean;default:1;" json:"downloaded"` // true when file is ready to be served, false when still being downloaded
}

// FullExternalUrl is the full url of the file after being downloaded for external use
func (f *File) FullExternalUrl() string {
	if !f.Downloaded.Valid || !f.Downloaded.Bool {
		return ""
	}
	return fmt.Sprintf("https://api.tum.app/files/%s%s", f.Path, f.Name)
}
