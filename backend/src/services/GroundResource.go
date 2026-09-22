package services

import (
	"gorm.io/gorm"

	"groundTurn/src/constants"
	"groundTurn/src/constructors"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
)

type GroundResourceService struct {
	db   *gorm.DB
	repo *repositories.GroundResourceRepository
}

func NewGroundResourceService(db *gorm.DB, repo *repositories.GroundResourceRepository) *GroundResourceService {
	return &GroundResourceService{db: db, repo: repo}
}

func (s *GroundResourceService) List() ([]types.GroundResourceDTO, error) {
	rows, err := s.repo.List()
	if err != nil {
		return nil, types.NewAppError(constants.InternalError, constants.InternalErrorMessage)
	}
	return constructors.NewGroundResourceListResponse(rows), nil
}
