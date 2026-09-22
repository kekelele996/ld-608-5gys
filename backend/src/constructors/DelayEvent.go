package constructors

import (
	"groundTurn/src/models"
	"groundTurn/src/types"
)

// NewDelayEventResponse 延误事件响应构造器。
func NewDelayEventResponse(delay models.DelayEvent) types.DelayEventDTO {
	return types.DelayEventDTO{
		ID:                 delay.ID,
		TurnaroundID:       delay.TurnaroundID,
		DelayType:          delay.DelayType,
		Minutes:            delay.Minutes,
		RootCause:          delay.RootCause,
		ResponsibilityTeam: delay.ResponsibilityTeam,
		ResolvedAt:         delay.ResolvedAt,
	}
}

func NewDelayEventListResponse(rows []models.DelayEvent) []types.DelayEventDTO {
	out := make([]types.DelayEventDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, NewDelayEventResponse(row))
	}
	return out
}
