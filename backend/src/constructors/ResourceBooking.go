package constructors

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
)

func NewBookingBlocker(booking models.ResourceBooking, resourceCode, reason string) types.BookingBlocker {
	return types.BookingBlocker{ResourceBooking: booking, ResourceCode: resourceCode, Reason: reason}
}

func ReleaseBooking(booking *models.ResourceBooking, releasedAt string) {
	booking.BookingStatus = constants.BookingReleased
	booking.ConflictReason = ""
	booking.ReleasedAt = &releasedAt
}
