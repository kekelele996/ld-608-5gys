package types

import "groundTurn/src/models"

type TaskBlocker struct {
	models.GroundTask
	Reason string `json:"reason"`
}
