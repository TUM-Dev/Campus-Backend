package model

import (
	"database/sql"
)

type Feedback struct {
	Id         int32           `gorm:"column:id;primary_key;AUTO_INCREMENT;type:int;not null;"`
	ImageCount int32           `gorm:"column:image_count;type:int;not null;"`
	EmailId    sql.NullString  `gorm:"column:email_id;type:text;null"`
	Receiver   sql.NullString  `gorm:"column:receiver;type:text;null"`
	ReplyTo    sql.NullString  `gorm:"column:reply_to;type:text;null"`
	Feedback   sql.NullString  `gorm:"column:feedback;type:text;null"`
	Latitude   sql.NullFloat64 `gorm:"column:latitude;type:float;null;"`
	Longitude  sql.NullFloat64 `gorm:"column:longitude;type:float;null;"`
	OsVersion  sql.NullString  `gorm:"column:os_version;type:text;null;"`
	AppVersion sql.NullString  `gorm:"column:app_version;type:text;null;"`
	Processed  bool            `gorm:"column:processed;type:boolean;default:false;not null;"`
	Timestamp  sql.NullTime    `gorm:"column:timestamp;type:timestamp;default:CURRENT_TIMESTAMP;null;"`
}

// TableName sets the insert table name for this struct type
func (n *Feedback) TableName() string {
	return "feedback"
}
