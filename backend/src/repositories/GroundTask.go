package repositories

import (
	"time"

	"gorm.io/gorm"

	"groundTurn/src/models"
)

type GroundTaskRepository struct{ db *gorm.DB }

func NewGroundTaskRepository(db *gorm.DB) *GroundTaskRepository {
	return &GroundTaskRepository{db: db}
}

func (r *GroundTaskRepository) List() ([]models.GroundTask, error) {
	var rows []models.GroundTask
	err := r.db.Order("turnaround_id asc, deadline asc").Find(&rows).Error
	return rows, err
}

func (r *GroundTaskRepository) ListByTurnaround(tx *gorm.DB, turnaroundID int64) ([]models.GroundTask, error) {
	var rows []models.GroundTask
	err := tx.Where("turnaround_id = ?", turnaroundID).Order("deadline asc").Find(&rows).Error
	return rows, err
}

func (r *GroundTaskRepository) GetByID(tx *gorm.DB, id int64) (*models.GroundTask, error) {
	var row models.GroundTask
	if err := tx.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// Sign 签收任务并写入 actual_finish；已是 SIGNED 的任务返回 0 行，避免重复签收。
func (r *GroundTaskRepository) Sign(tx *gorm.DB, id int64, finishTime time.Time) (int64, error) {
	result := tx.Model(&models.GroundTask{}).
		Where("id = ? AND status <> ?", id, "SIGNED").
		Updates(map[string]interface{}{
			"status":        "SIGNED",
			"actual_finish": finishTime,
		})
	return result.RowsAffected, result.Error
}

// TouchOnRelease 放行事务里同步刷新任务时间线（updated_at 与过站 READY 对齐）。
func (r *GroundTaskRepository) TouchOnRelease(tx *gorm.DB, turnaroundID int64, at time.Time) error {
	return tx.Model(&models.GroundTask{}).
		Where("turnaround_id = ?", turnaroundID).
		Update("updated_at", at).Error
}
