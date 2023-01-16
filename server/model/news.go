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
	//[ 0] news                                           int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	News int32 `gorm:"primary_key;AUTO_INCREMENT;column:news;type:int;" json:"news"`
	//[ 1] date                                           datetime             null: false  primary: false  isArray: false  auto: false  col: datetime        len: -1      default: []
	Date time.Time `gorm:"column:date;type:datetime;" json:"date"`
	//[ 2] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 3] title                                          text(255)            null: false  primary: false  isArray: false  auto: false  col: text            len: 255     default: []
	Title string `gorm:"column:title;type:text;size:255;" json:"title"`
	//[ 4] description                                    text(65535)          null: false  primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Description string `gorm:"column:description;type:text;size:65535;" json:"description"`
	//[ 5] src                                            int                  null: false  primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	Src int32 `gorm:"column:src;type:int;" json:"src"`
	//[ 6] link                                           varchar(190)         null: false  primary: false  isArray: false  auto: false  col: varchar         len: 190     default: []
	Link string `gorm:"column:link;type:varchar(190);" json:"link"`
	//[ 7] image                                          text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Image null.String `gorm:"column:image;type:text;size:65535;" json:"image"`
	//[ 8] file                                           int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	File null.Int `gorm:"column:file;type:int;" json:"file"`
}

// TableName sets the insert table name for this struct type
func (n *News) TableName() string {
	return "news"
}
