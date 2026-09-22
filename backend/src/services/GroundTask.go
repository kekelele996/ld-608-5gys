package services

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
)

type GroundTaskService struct{ store *repositories.Store }

func NewGroundTaskService(store *repositories.Store) *GroundTaskService {
	return &GroundTaskService{store: store}
}

func (s *GroundTaskService) List() []models.GroundTask { return s.store.ListTasks() }

func (s *GroundTaskService) UnsignedCount(turnaroundID int) int {
	count := 0
	for _, task := range s.store.ListTasksByTurnaround(turnaroundID) {
		if task.Status != constants.GroundTaskSigned {
			count++
		}
	}
	return count
}
