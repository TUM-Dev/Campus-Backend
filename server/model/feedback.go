package model

import (
	"github.com/guregu/null"
)

type Feedback struct {
	Id           int64       `gorm:"column:id;primary_key;AUTO_INCREMENT;type:int;not null;"`
	ImageCount   int32       `gorm:"column:image_count;type:int;not null;"`
	EmailId      string      `gorm:"column:email_id;type:text;not null"`
	Recipient    string      `gorm:"column:receiver;type:text;not null;uniqueIndex:receiver_reply_to_feedback_app_version_uindex"`
	ReplyToEmail null.String `gorm:"column:reply_to_email;type:text;null;uniqueIndex:receiver_reply_to_feedback_app_version_uindex"`
	ReplyToName  null.String `gorm:"column:reply_to_name;type:text;null"`
	Feedback     string      `gorm:"column:feedback;type:text;not null;uniqueIndex:receiver_reply_to_feedback_app_version_uindex"`
	Latitude     null.Float  `gorm:"column:latitude;type:float;null;"`
	Longitude    null.Float  `gorm:"column:longitude;type:float;null;"`
	OsVersion    null.String `gorm:"column:os_version;type:text;null;"`
	AppVersion   null.String `gorm:"column:app_version;type:text;null;uniqueIndex:receiver_reply_to_feedback_app_version_uindex"`
	Processed    bool        `gorm:"column:processed;type:boolean;default:false;not null;"`
	Timestamp    null.Time   `gorm:"column:timestamp;type:timestamp;default:CURRENT_TIMESTAMP;null;"`
}

// TableName sets the insert table name for this struct type
func (n *Feedback) TableName() string {
	return "feedback"
}
