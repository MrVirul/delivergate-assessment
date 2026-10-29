package models

import "time"

type InventoryItem struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"type:varchar(120);not null;uniqueIndex" json:"name"`
	Category     string    `gorm:"type:varchar(30);not null;index" json:"category"`
	Description  string    `gorm:"type:text" json:"description"`
	Quantity     int       `gorm:"not null" json:"quantity"`
	Unit         string    `gorm:"type:varchar(20);not null" json:"unit"`
	PriceCents   int64     `gorm:"not null" json:"price_cents"`
	ReorderLevel int       `gorm:"not null;default:0" json:"reorder_level"`
	CreatedByID  *uint     `json:"created_by_id"`
	UpdatedByID  *uint     `json:"updated_by_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
