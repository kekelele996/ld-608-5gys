package constructors

import (
	"groundTurn/src/models"
	"groundTurn/src/types"
)

func NewDelayBlocker(delay models.DelayEvent, reason string) types.DelayBlocker {
	return types.DelayBlocker{DelayEvent: delay, Reason: reason}
}
