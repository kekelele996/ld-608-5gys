package models

type GroundTask struct {
	ID               int     `json:"id"`
	TurnaroundID     int     `json:"turnaround_id"`
	TaskType         string  `json:"task_type"`
	TeamID           int     `json:"team_id"`
	PlannedStart     string  `json:"planned_start"`
	Deadline         string  `json:"deadline"`
	ActualFinish     *string `json:"actual_finish,omitempty"`
	Status           string  `json:"status"`
	BlockerNote      string  `json:"blocker_note"`
	TimelineSyncedAt *string `json:"timeline_synced_at,omitempty"`
}
