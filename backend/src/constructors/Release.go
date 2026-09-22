package constructors

import (
	"groundTurn/src/models"
	"groundTurn/src/types"
)

func NewReleaseResponse(flight types.FlightTurnaroundDetail, tasks []models.GroundTask, bookings []models.ResourceBooking, releasedCount int) types.ReleaseResponse {
	return types.ReleaseResponse{
		Message:       "放行成功，航班进入 READY",
		Turnaround:    flight,
		Tasks:         tasks,
		Bookings:      bookings,
		ReleasedCount: releasedCount,
	}
}
