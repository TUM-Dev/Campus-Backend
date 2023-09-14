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

// Files struct is a row record of the files table in the tca database
type Files struct {
	//[ 0] file                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	File int32 `gorm:"primary_key;AUTO_INCREMENT;column:file;type:int;" json:"file"`
	//[ 1] name                                           text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Name string `gorm:"column:name;type:text;size:16777215;" json:"name"`
	//[ 2] path                                           text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Path string `gorm:"column:path;type:text;size:16777215;" json:"path"`
	//[ 3] downloads                                      int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: [0]
	Downloads  int32          `gorm:"column:downloads;type:int;default:0;" json:"downloads"`
	URL        sql.NullString `gorm:"column:url;default:null;" json:"url"`                         // URL of the files source (if any)
	Downloaded sql.NullBool   `gorm:"column:downloaded;type:boolean;default:1;" json:"downloaded"` // true when file is ready to be served, false when still being downloaded
}

// TableName sets the insert table name for this struct type
func (f *Files) TableName() string {
	return "files"
}
