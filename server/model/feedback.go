package model

import (
	"github.com/guregu/null"
)

type Feedback struct {
	Id           int64       `gorm:"column:id;primary_key;AUTO_INCREMENT;not null;"`
	ImageCount   int32       `gorm:"column:image_count;not null;"`
	EmailId      string      `gorm:"column:email_id;not null"`
	Recipient    string      `gorm:"column:receiver;not null;uniqueIndex:receiver_reply_to_feedback_app_version_uindex,expression:receiver(255)"`
	ReplyToEmail null.String `gorm:"column:reply_to_email;null;uniqueIndex:receiver_reply_to_feedback_app_version_uindex,expression:reply_to_email(100)"`
	ReplyToName  null.String `gorm:"column:reply_to_name;null"`
	Feedback     string      `gorm:"column:feedback;not null;uniqueIndex:receiver_reply_to_feedback_app_version_uindex,expression:feedback(255)"`
	Latitude     null.Float  `gorm:"column:latitude;type:double;null;"`
	Longitude    null.Float  `gorm:"column:longitude;type:double;null;"`
	OsVersion    null.String `gorm:"column:os_version;null;"`
	AppVersion   null.String `gorm:"column:app_version;null;uniqueIndex:receiver_reply_to_feedback_app_version_uindex,expression:app_version(100)"`
	Processed    bool        `gorm:"column:processed;default:false;not null;"`
	Timestamp    null.Time   `gorm:"column:timestamp;type:datetime;default:current_timestamp();null;"`
}

// TableName sets the insert table name for this struct type
func (n *Feedback) TableName() string {
	return "feedback"
}
