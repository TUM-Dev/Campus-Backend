package model

import "time"

type AdminSession struct {
	TokenHash string    `gorm:"column:token_hash;type:char(64);primaryKey"`
	Issuer    string    `gorm:"column:issuer;type:varchar(255);not null"`
	Subject   string    `gorm:"column:subject;type:varchar(255);not null"`
	CSRFToken string    `gorm:"column:csrf_token;type:char(64);not null"`
	ExpiresAt time.Time `gorm:"column:expires_at;type:datetime(6);not null;index"`
}
