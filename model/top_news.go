package model

import (
	"time"
)

type TopNews struct {
	Name    string
	Link    string `protobuf:"bytes,2,opt,name=link,proto3" json:"link,omitempty"`
	Created *time.Time
	From    *time.Time
	To      *time.Time
}
