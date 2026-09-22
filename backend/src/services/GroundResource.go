package services

import (
	"groundTurn/src/repositories"
	"groundTurn/src/types"
)

type GroundResourceService struct{ store *repositories.Store }

func NewGroundResourceService(store *repositories.Store) *GroundResourceService {
	return &GroundResourceService{store: store}
}

func (s *GroundResourceService) List() []types.ResourceView {
	snapshot := s.store.Snapshot()
	activeByResource := map[int]int{}
	for _, booking := range snapshot.Bookings {
		if booking.BookingStatus == "ACTIVE" {
			activeByResource[booking.ResourceID]++
		}
	}
	views := make([]types.ResourceView, 0, len(snapshot.Resources))
	for _, resource := range snapshot.Resources {
		views = append(views, types.ResourceView{GroundResource: resource, ActiveBookingCount: activeByResource[resource.ID]})
	}
	return views
}
