package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	"github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

// NewsSource struct is a row record of the newsSource table in the tca database
type NewsSource struct {
	//[ 0] source                                         int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	Source int32 `gorm:"primary_key;AUTO_INCREMENT;column:source;type:int;" json:"source"`
	//[ 1] title                                          text(16777215)       null: false  primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	Title string `gorm:"column:title;type:text;size:16777215;" json:"title"`
	//[ 2] url                                            text(16777215)       null: true   primary: false  isArray: false  auto: false  col: text            len: 16777215 default: []
	URL null.String `gorm:"column:url;type:text;size:16777215;" json:"url"`
	//[ 3] icon                                           int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Icon null.Int `gorm:"column:icon;type:int;" json:"icon"`
	//[ 4] hook                                           char(12)             null: true   primary: false  isArray: false  auto: false  col: char            len: 12      default: []
	Hook null.String `gorm:"column:hook;type:char;size:12;" json:"hook"`
}

// TableName sets the insert table name for this struct type
func (n *NewsSource) TableName() string {
	return "newsSource"
}
