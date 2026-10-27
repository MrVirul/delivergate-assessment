package repository

import (
	"time"

	"github.com/MrVirul/icecream-inventory/internal/models"
	"gorm.io/gorm"
)

type RevokedTokenRepository interface {
	Create(jti string, expiresAt time.Time) error
	IsRevoked(jti string) (bool, error)
	DeleteExpired() error
}

type revokedTokenRepository struct {
	db *gorm.DB
}

func NewRevokedTokenRepository(db *gorm.DB) RevokedTokenRepository {
	return &revokedTokenRepository{db: db}
}

func (r *revokedTokenRepository) Create(jti string, expiresAt time.Time) error {
	rt := &models.RevokedToken{JTI: jti, ExpiresAt: expiresAt}
	// idempotent
	return r.db.FirstOrCreate(rt, models.RevokedToken{JTI: jti}).Error
}

func (r *revokedTokenRepository) IsRevoked(jti string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.RevokedToken{}).Where("jti = ?", jti).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *revokedTokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.RevokedToken{}).Error
}
