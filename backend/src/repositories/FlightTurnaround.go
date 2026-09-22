package repositories

import "groundTurn/src/models"

func (s *Store) ListFlights() []models.FlightTurnaround {
	return s.Snapshot().Flights
}

func (s *Store) GetFlight(id int) (models.FlightTurnaround, bool) {
	for _, flight := range s.Snapshot().Flights {
		if flight.ID == id {
			return flight, true
		}
	}
	return models.FlightTurnaround{}, false
}
