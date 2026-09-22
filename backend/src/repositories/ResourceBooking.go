package repositories

import (
	"time"

	"gorm.io/gorm"

	"groundTurn/src/models"
)

type ResourceBookingRepository struct{ db *gorm.DB }

func NewResourceBookingRepository(db *gorm.DB) *ResourceBookingRepository {
	return &ResourceBookingRepository{db: db}
}

func (r *ResourceBookingRepository) List() ([]models.ResourceBooking, error) {
	var rows []models.ResourceBooking
	err := r.db.Preload("Resource").Order("turnaround_id asc, start_time asc").Find(&rows).Error
	return rows, err
}

func (r *ResourceBookingRepository) GetByID(tx *gorm.DB, id int64) (*models.ResourceBooking, error) {
	var row models.ResourceBooking
	if err := tx.Preload("Resource").First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// ResolveConflict 人工解决冲突：CONFLICT -> ACTIVE（之后仍由放行事务一次性释放）。
func (r *ResourceBookingRepository) ResolveConflict(tx *gorm.DB, id int64) (int64, error) {
	result := tx.Model(&models.ResourceBooking{}).
		Where("id = ? AND booking_status = ?", id, "CONFLICT").
		Updates(map[string]interface{}{
			"booking_status":  "ACTIVE",
			"conflict_reason": "",
		})
	return result.RowsAffected, result.Error
}

// ReleaseActiveForTurnaround 一次性释放航班下所有 ACTIVE 预约（CONFLICT 不自动释放），返回释放数量与预约 ID。
func (r *ResourceBookingRepository) ReleaseActiveForTurnaround(tx *gorm.DB, turnaroundID int64, at time.Time) ([]int64, error) {
	var ids []int64
	if err := tx.Model(&models.ResourceBooking{}).
		Where("turnaround_id = ? AND booking_status = ?", turnaroundID, "ACTIVE").
		Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, nil
	}
	if err := tx.Model(&models.ResourceBooking{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"booking_status": "RELEASED",
			"released_at":    at,
		}).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// MarkResourcesAvailable 释放预约后把对应资源台账刷回 AVAILABLE。
func (r *ResourceBookingRepository) MarkResourcesAvailable(tx *gorm.DB, resourceIDs []int64) error {
	if len(resourceIDs) == 0 {
		return nil
	}
	return tx.Model(&models.GroundResource{}).
		Where("id IN ?", resourceIDs).
		Update("availability_status", "AVAILABLE").Error
}

func (r *ResourceBookingRepository) ListResourcesByBookingIDs(tx *gorm.DB, ids []int64) ([]int64, error) {
	var resourceIDs []int64
	if len(ids) == 0 {
		return resourceIDs, nil
	}
	err := tx.Model(&models.ResourceBooking{}).
		Where("id IN ?", ids).
		Distinct().
		Pluck("resource_id", &resourceIDs).Error
	return resourceIDs, err
}
