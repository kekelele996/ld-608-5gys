package repositories

import (
	"gorm.io/gorm"

	"groundTurn/src/models"
)

type FlightTurnaroundRepository struct{ db *gorm.DB }

func NewFlightTurnaroundRepository(db *gorm.DB) *FlightTurnaroundRepository {
	return &FlightTurnaroundRepository{db: db}
}

func (r *FlightTurnaroundRepository) List() ([]models.FlightTurnaround, error) {
	var rows []models.FlightTurnaround
	err := r.db.Order("arrival_time asc, id asc").Find(&rows).Error
	return rows, err
}

// GetWithAssociations 拉取航班下全部任务、预约与延误，供放行门禁聚合。
func (r *FlightTurnaroundRepository) GetWithAssociations(id int64) (*models.FlightTurnaround, error) {
	var row models.FlightTurnaround
	err := r.db.
		Preload("Tasks").
		Preload("Bookings").
		Preload("Bookings.Resource").
		Preload("Delays").
		First(&row, id).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *FlightTurnaroundRepository) GetByID(tx *gorm.DB, id int64) (*models.FlightTurnaround, error) {
	var row models.FlightTurnaround
	err := tx.First(&row, id).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// MarkReady 条件更新：仅当前状态不在 READY/DEPARTED 时生效，保证并发/重复放行只生效一次。返回受影响行数。
func (r *FlightTurnaroundRepository) MarkReady(tx *gorm.DB, id int64, readyAt interface{}) (int64, error) {
	result := tx.Model(&models.FlightTurnaround{}).
		Where("id = ? AND turnaround_status <> ? AND turnaround_status <> ?",
			id, "READY", "DEPARTED").
		Updates(map[string]interface{}{
			"turnaround_status": "READY",
			"ready_at":          readyAt,
		})
	return result.RowsAffected, result.Error
}
