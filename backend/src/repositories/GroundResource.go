package repositories

import (
	"gorm.io/gorm"

	"groundTurn/src/models"
)

type GroundResourceRepository struct{ db *gorm.DB }

func NewGroundResourceRepository(db *gorm.DB) *GroundResourceRepository {
	return &GroundResourceRepository{db: db}
}

func (r *GroundResourceRepository) List() ([]models.GroundResource, error) {
	var rows []models.GroundResource
	err := r.db.Order("resource_code asc").Find(&rows).Error
	return rows, err
}

func (r *GroundResourceRepository) CreateBatch(tx *gorm.DB, rows []models.GroundResource) error {
	if len(rows) == 0 {
		return nil
	}
	return tx.Create(&rows).Error
}
