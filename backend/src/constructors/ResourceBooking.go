package constructors

import (
	"groundTurn/src/models"
	"groundTurn/src/types"
)

// NewResourceBookingResponse 资源预约响应构造器。
func NewResourceBookingResponse(booking models.ResourceBooking) types.ResourceBookingDTO {
	resourceCode := ""
	if booking.Resource != nil {
		resourceCode = booking.Resource.ResourceCode
	}
	return types.ResourceBookingDTO{
		ID:             booking.ID,
		ResourceID:     booking.ResourceID,
		ResourceCode:   resourceCode,
		TurnaroundID:   booking.TurnaroundID,
		TaskID:         booking.TaskID,
		StartTime:      booking.StartTime,
		EndTime:        booking.EndTime,
		BookingStatus:  booking.BookingStatus,
		ConflictReason: booking.ConflictReason,
		ReleasedAt:     booking.ReleasedAt,
	}
}

func NewResourceBookingListResponse(rows []models.ResourceBooking) []types.ResourceBookingDTO {
	out := make([]types.ResourceBookingDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, NewResourceBookingResponse(row))
	}
	return out
}
