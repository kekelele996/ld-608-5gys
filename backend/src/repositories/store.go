package repositories

import (
	"sync"

	"groundTurn/src/models"
)

type Data struct {
	Flights   []models.FlightTurnaround
	Tasks     []models.GroundTask
	Resources []models.GroundResource
	Bookings  []models.ResourceBooking
	Delays    []models.DelayEvent
	Logs      []models.AuditLog
}

type Store struct {
	mu   sync.RWMutex
	data Data
}

func NewStore(data Data) *Store {
	return &Store{data: cloneData(data)}
}

func (s *Store) Snapshot() Data {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneData(s.data)
}

func (s *Store) Update(mutate func(*Data) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	working := cloneData(s.data)
	if err := mutate(&working); err != nil {
		return err
	}
	s.data = working
	return nil
}

func StringPtr(value string) *string { return &value }

func cloneData(source Data) Data {
	clone := Data{
		Flights:   append([]models.FlightTurnaround(nil), source.Flights...),
		Tasks:     append([]models.GroundTask(nil), source.Tasks...),
		Resources: append([]models.GroundResource(nil), source.Resources...),
		Bookings:  append([]models.ResourceBooking(nil), source.Bookings...),
		Delays:    append([]models.DelayEvent(nil), source.Delays...),
		Logs:      append([]models.AuditLog(nil), source.Logs...),
	}
	for i := range clone.Flights {
		clone.Flights[i].ReadyAt = cloneStringPtr(source.Flights[i].ReadyAt)
		clone.Flights[i].TimelineSyncedAt = cloneStringPtr(source.Flights[i].TimelineSyncedAt)
	}
	for i := range clone.Tasks {
		clone.Tasks[i].ActualFinish = cloneStringPtr(source.Tasks[i].ActualFinish)
		clone.Tasks[i].TimelineSyncedAt = cloneStringPtr(source.Tasks[i].TimelineSyncedAt)
	}
	for i := range clone.Bookings {
		clone.Bookings[i].ReleasedAt = cloneStringPtr(source.Bookings[i].ReleasedAt)
	}
	for i := range clone.Delays {
		clone.Delays[i].ResolvedAt = cloneStringPtr(source.Delays[i].ResolvedAt)
	}
	return clone
}

func cloneStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}
