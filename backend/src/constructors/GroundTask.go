package constructors

import (
	"groundTurn/src/models"
	"groundTurn/src/types"
)

// NewGroundTaskResponse 地勤任务响应构造器，页面/store/service 禁止散写 DTO。
func NewGroundTaskResponse(task models.GroundTask) types.GroundTaskDTO {
	return types.GroundTaskDTO{
		ID:           task.ID,
		TurnaroundID: task.TurnaroundID,
		TaskType:     task.TaskType,
		TeamID:       task.TeamID,
		PlannedStart: task.PlannedStart,
		Deadline:     task.Deadline,
		ActualFinish: task.ActualFinish,
		Status:       task.Status,
		BlockerNote:  task.BlockerNote,
	}
}

func NewGroundTaskListResponse(rows []models.GroundTask) []types.GroundTaskDTO {
	out := make([]types.GroundTaskDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, NewGroundTaskResponse(row))
	}
	return out
}
