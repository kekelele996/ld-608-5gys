package types

import "groundTurn/src/models"

type BookingBlocker struct {
	models.ResourceBooking
	ResourceCode string `json:"resource_code"`
	Reason       string `json:"reason"`
}
