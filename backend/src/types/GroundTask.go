package types

import "time"

type GroundTaskDTO struct {
	ID           int64      `json:"id"`
	TurnaroundID int64      `json:"turnaround_id"`
	TaskType     string     `json:"task_type"`
	TeamID       int64      `json:"team_id"`
	PlannedStart time.Time  `json:"planned_start"`
	Deadline     time.Time  `json:"deadline"`
	ActualFinish *time.Time `json:"actual_finish"`
	Status       string     `json:"status"`
	BlockerNote  string     `json:"blocker_note"`
}

// SignTaskRequest 签收完成任务。SIGNED 是放行门禁的唯一合法任务终态。
type SignTaskRequest struct {
	Actor string `json:"actor"`
}
