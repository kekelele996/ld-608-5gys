package types

import "groundTurn/src/models"

type ResourceView struct {
	models.GroundResource
	ActiveBookingCount int `json:"active_booking_count"`
}
