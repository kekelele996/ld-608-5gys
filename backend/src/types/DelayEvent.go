package types

import "groundTurn/src/models"

type DelayBlocker struct {
	models.DelayEvent
	Reason string `json:"reason"`
}
