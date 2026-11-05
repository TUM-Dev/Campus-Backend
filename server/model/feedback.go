package model

import (
	"time"

	"github.com/guregu/null"
)

type Feedback struct {
	Id           int64       `gorm:"column:id;primary_key;autoIncrement;not null"`
	ImageCount   int32       `gorm:"column:image_count;not null"`
	EmailId      string      `gorm:"column:email_id;type:text"`
	Recipient    string      `gorm:"column:receiver;type:text;uniqueIndex:receiver_reply_to_feedback_app_version_uindex,expression:receiver(255)"`
	ReplyToEmail null.String `gorm:"column:reply_to_email;type:text;uniqueIndex:receiver_reply_to_feedback_app_version_uindex,expression:reply_to_email(100)"`
	ReplyToName  null.String `gorm:"column:reply_to_name;type:text"`
	Feedback     string      `gorm:"column:feedback;type:text;uniqueIndex:receiver_reply_to_feedback_app_version_uindex,expression:feedback(255)"`
	Latitude     null.Float  `gorm:"column:latitude;type:double"`
	Longitude    null.Float  `gorm:"column:longitude;type:double"`
	OsVersion    null.String `gorm:"column:os_version;type:text"`
	AppVersion   null.String `gorm:"column:app_version;type:text;uniqueIndex:receiver_reply_to_feedback_app_version_uindex,expression:app_version(100)"`
	Processed    bool        `gorm:"column:processed;default:false;not null"`
	Timestamp    time.Time   `gorm:"column:timestamp;type:datetime;default:current_timestamp();not null"`
}

// TableName sets the insert table name for this struct type
func (n *Feedback) TableName() string {
	return "feedback"
}
