package model

import (
	"time"

	"github.com/guregu/null"
)

type Feedback struct {
	Id           int64       `gorm:"column:id;primary_key;autoIncrement;not null"`
	ImageCount   int32       `gorm:"column:image_count;not null"`
	EmailId      string      `gorm:"column:email_id;not null"`
	Recipient    string      `gorm:"column:receiver;not null;uniqueIndex:receiver_reply_to_feedback_app_version_uindex,expression:receiver(255)"`
	ReplyToEmail null.String `gorm:"column:reply_to_email;uniqueIndex:receiver_reply_to_feedback_app_version_uindex,expression:reply_to_email(100)"`
	ReplyToName  null.String `gorm:"column:reply_to_name"`
	Feedback     string      `gorm:"column:feedback;not null;uniqueIndex:receiver_reply_to_feedback_app_version_uindex,expression:feedback(255)"`
	Latitude     null.Float  `gorm:"column:latitude;type:double"`
	Longitude    null.Float  `gorm:"column:longitude;type:double"`
	OsVersion    null.String `gorm:"column:os_version"`
	AppVersion   null.String `gorm:"column:app_version;uniqueIndex:receiver_reply_to_feedback_app_version_uindex,expression:app_version(100)"`
	Processed    bool        `gorm:"column:processed;default:false;not null"`
	Timestamp    time.Time   `gorm:"column:timestamp;type:datetime;default:current_timestamp();not null"`
}

// TableName sets the insert table name for this struct type
func (n *Feedback) TableName() string {
	return "feedback"
}
