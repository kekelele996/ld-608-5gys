package types

import (
	"net/http"

	"groundTurn/src/models"
)

type ReleaseBlockers struct {
	Tasks    []TaskBlocker    `json:"tasks"`
	Delays   []DelayBlocker   `json:"delays"`
	Bookings []BookingBlocker `json:"bookings"`
}

type ReleaseResponse struct {
	Message       string                   `json:"message"`
	Turnaround    FlightTurnaroundDetail   `json:"turnaround"`
	Tasks         []models.GroundTask      `json:"tasks"`
	Bookings      []models.ResourceBooking `json:"bookings"`
	ReleasedCount int                      `json:"released_count"`
}

type ReleaseRejectedResponse struct {
	Code     string          `json:"code"`
	Message  string          `json:"message"`
	Blockers ReleaseBlockers `json:"blockers"`
	HTTPCode int             `json:"-"`
}

func NewReleaseRejected(message string, blockers ReleaseBlockers) *ReleaseRejectedResponse {
	return &ReleaseRejectedResponse{Code: "RELEASE_BLOCKED", Message: message, Blockers: blockers, HTTPCode: http.StatusConflict}
}
