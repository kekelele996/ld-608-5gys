package services

import (
	"groundTurn/src/models"
	"groundTurn/src/repositories"
)

type ResourceBookingService struct{ store *repositories.Store }

func NewResourceBookingService(store *repositories.Store) *ResourceBookingService {
	return &ResourceBookingService{store: store}
}

func (s *ResourceBookingService) List() []models.ResourceBooking { return s.store.ListBookings() }
