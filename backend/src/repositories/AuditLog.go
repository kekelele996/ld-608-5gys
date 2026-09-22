package repositories

import (
	"gorm.io/gorm"

	"groundTurn/src/models"
)

// AuditLogRepository 写操作审计落库。
type AuditLogRepository struct{ db *gorm.DB }

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Append(tx *gorm.DB, log models.AuditLog) error {
	return tx.Create(&log).Error
}
