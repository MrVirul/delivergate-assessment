package models

import "time"

type RevokedToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	JTI       string    `gorm:"type:varchar(64);not null;uniqueIndex" json:"jti"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
