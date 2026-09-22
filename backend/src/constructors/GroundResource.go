package constructors

import (
	"groundTurn/src/models"
	"groundTurn/src/types"
)

func NewResourceView(resource models.GroundResource, activeBookingCount int) types.ResourceView {
	return types.ResourceView{GroundResource: resource, ActiveBookingCount: activeBookingCount}
}

func ResourceCodeMap(resources []models.GroundResource) map[int]string {
	names := make(map[int]string, len(resources))
	for _, resource := range resources {
		names[resource.ID] = resource.ResourceCode
	}
	return names
}
