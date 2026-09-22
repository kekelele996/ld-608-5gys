package repositories

import (
	"time"

	"gorm.io/gorm"

	"groundTurn/src/models"
)

type DelayEventRepository struct{ db *gorm.DB }

func NewDelayEventRepository(db *gorm.DB) *DelayEventRepository {
	return &DelayEventRepository{db: db}
}

func (r *DelayEventRepository) List() ([]models.DelayEvent, error) {
	var rows []models.DelayEvent
	err := r.db.Order("turnaround_id asc, created_at asc").Find(&rows).Error
	return rows, err
}

func (r *DelayEventRepository) GetByID(tx *gorm.DB, id int64) (*models.DelayEvent, error) {
	var row models.DelayEvent
	if err := tx.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Resolve 关闭延误：仅对 resolved_at 为空的未关闭事件生效。
func (r *DelayEventRepository) Resolve(tx *gorm.DB, id int64, at time.Time) (int64, error) {
	result := tx.Model(&models.DelayEvent{}).
		Where("id = ? AND resolved_at IS NULL", id).
		Update("resolved_at", at)
	return result.RowsAffected, result.Error
}
