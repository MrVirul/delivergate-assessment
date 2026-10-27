package models

import "time"

type Customer struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(120);not null" json:"name"`
	Email        *string   `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	Phone        string    `gorm:"type:varchar(30);not null;index" json:"phone"`
	Address      string    `gorm:"type:varchar(255)" json:"address"`
	Notes        string    `gorm:"type:text" json:"notes"`
	CreatedByID  *uint     `json:"created_by_id"`
	UpdatedByID  *uint     `json:"updated_by_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
