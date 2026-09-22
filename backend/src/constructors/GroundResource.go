package constructors

import (
	"groundTurn/src/models"
	"groundTurn/src/types"
)

// NewGroundResourceResponse 保障资源响应构造器。
func NewGroundResourceResponse(resource models.GroundResource) types.GroundResourceDTO {
	return types.GroundResourceDTO{
		ID:                 resource.ID,
		ResourceCode:       resource.ResourceCode,
		ResourceType:       resource.ResourceType,
		Location:           resource.Location,
		AvailabilityStatus: resource.AvailabilityStatus,
		MaintenanceDueAt:   resource.MaintenanceDueAt,
		OwnerTeam:          resource.OwnerTeam,
	}
}

func NewGroundResourceListResponse(rows []models.GroundResource) []types.GroundResourceDTO {
	out := make([]types.GroundResourceDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, NewGroundResourceResponse(row))
	}
	return out
}
