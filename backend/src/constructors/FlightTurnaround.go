package constructors

import (
	"groundTurn/src/models"
	"groundTurn/src/types"
)

func NewFlightDetail(flight models.FlightTurnaround, timeline []types.TimelineEvent, taskCompletionRate int) types.FlightTurnaroundDetail {
	return types.FlightTurnaroundDetail{
		FlightTurnaround:   flight,
		TaskCompletionRate: taskCompletionRate,
		Timeline:           timeline,
	}
}
