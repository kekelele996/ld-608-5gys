package repositories

import "groundTurn/src/models"

func (s *Store) ListResources() []models.GroundResource {
	return s.Snapshot().Resources
}

func (s *Store) GetResource(id int) (models.GroundResource, bool) {
	for _, resource := range s.Snapshot().Resources {
		if resource.ID == id {
			return resource, true
		}
	}
	return models.GroundResource{}, false
}
