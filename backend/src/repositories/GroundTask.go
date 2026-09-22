package repositories

import "groundTurn/src/models"

func (s *Store) ListTasks() []models.GroundTask {
	return s.Snapshot().Tasks
}

func (s *Store) ListTasksByTurnaround(turnaroundID int) []models.GroundTask {
	var rows []models.GroundTask
	for _, task := range s.Snapshot().Tasks {
		if task.TurnaroundID == turnaroundID {
			rows = append(rows, task)
		}
	}
	return rows
}
