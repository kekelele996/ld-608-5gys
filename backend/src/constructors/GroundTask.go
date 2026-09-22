package constructors

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
)

func NewTaskBlocker(task models.GroundTask, reason string) types.TaskBlocker {
	return types.TaskBlocker{GroundTask: task, Reason: reason}
}

func SignTaskForRelease(task *models.GroundTask, syncedAt string) {
	task.Status = constants.GroundTaskSigned
	task.TimelineSyncedAt = &syncedAt
}
