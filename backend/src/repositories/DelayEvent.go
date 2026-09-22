package repositories

import "groundTurn/src/models"

func (s *Store) ListDelays() []models.DelayEvent {
	return s.Snapshot().Delays
}

func (s *Store) ListDelaysByTurnaround(turnaroundID int) []models.DelayEvent {
	var rows []models.DelayEvent
	for _, delay := range s.Snapshot().Delays {
		if delay.TurnaroundID == turnaroundID {
			rows = append(rows, delay)
		}
	}
	return rows
}
