package services

import (
	"groundTurn/src/models"
	"groundTurn/src/repositories"
)

type DelayEventService struct{ store *repositories.Store }

func NewDelayEventService(store *repositories.Store) *DelayEventService {
	return &DelayEventService{store: store}
}

func (s *DelayEventService) List() []models.DelayEvent { return s.store.ListDelays() }
