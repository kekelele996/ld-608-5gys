package constants

const (
	BookingActive   = "ACTIVE"
	BookingReleased = "RELEASED"
	BookingConflict = "CONFLICT"
)

var BookingStatuses = []string{BookingActive, BookingReleased, BookingConflict}
