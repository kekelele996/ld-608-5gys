package repositories

import "groundTurn/src/models"

func (s *Store) ListBookings() []models.ResourceBooking {
	return s.Snapshot().Bookings
}

func (s *Store) ListBookingsByTurnaround(turnaroundID int) []models.ResourceBooking {
	var rows []models.ResourceBooking
	for _, booking := range s.Snapshot().Bookings {
		if booking.TurnaroundID == turnaroundID {
			rows = append(rows, booking)
		}
	}
	return rows
}
