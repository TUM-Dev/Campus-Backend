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

// NewsAlert struct is a row record of the news_alert table in the tca database
type NewsAlert struct {
	//[ 0] news_alert                                     int                  null: false  primary: true   isArray: false  auto: true   col: int             len: -1      default: []
	NewsAlert int32 `gorm:"primary_key;AUTO_INCREMENT;column:news_alert;type:int;" json:"news_alert"`
	//[ 1] file                                           int                  null: true   primary: false  isArray: false  auto: false  col: int             len: -1      default: []
	File null.Int `gorm:"column:file;type:int;" json:"file"`
	//[ 2] name                                           varchar(100)         null: true   primary: false  isArray: false  auto: false  col: varchar         len: 100     default: []
	Name null.String `gorm:"column:name;type:varchar(100);" json:"name"`
	//[ 3] link                                           text(65535)          null: true   primary: false  isArray: false  auto: false  col: text            len: 65535   default: []
	Link null.String `gorm:"column:link;type:text;size:65535;" json:"link"`
	//[ 4] created                                        timestamp            null: false  primary: false  isArray: false  auto: false  col: timestamp       len: -1      default: [CURRENT_TIMESTAMP]
	Created time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;" json:"created"`
	//[ 5] from                                           datetime             null: false  primary: false  isArray: false  auto: false  col: datetime        len: -1      default: [CURRENT_TIMESTAMP]
	From time.Time `gorm:"column:from;type:datetime;default:CURRENT_TIMESTAMP;" json:"from"`
	//[ 6] to                                             datetime             null: false  primary: false  isArray: false  auto: false  col: datetime        len: -1      default: [CURRENT_TIMESTAMP]
	To time.Time `gorm:"column:to;type:datetime;default:CURRENT_TIMESTAMP;" json:"to"`
}

// TableName sets the insert table name for this struct type
func (n *NewsAlert) TableName() string {
	return "news_alert"
}
