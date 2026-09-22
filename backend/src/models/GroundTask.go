package models

import "time"

// GroundTask 地勤任务。status 取值见 constants.GroundTaskStatuses；放行前必须全部 SIGNED。
type GroundTask struct {
	ID           int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	TurnaroundID int64      `gorm:"column:turnaround_id;index" json:"turnaround_id"`
	TaskType     string     `gorm:"column:task_type;size:32" json:"task_type"`
	TeamID       int64      `gorm:"column:team_id" json:"team_id"`
	PlannedStart time.Time  `gorm:"column:planned_start" json:"planned_start"`
	Deadline     time.Time  `gorm:"column:deadline" json:"deadline"`
	ActualFinish *time.Time `gorm:"column:actual_finish" json:"actual_finish"`
	Status       string     `gorm:"column:status;size:32;index" json:"status"`
	BlockerNote  string     `gorm:"column:blocker_note;size:255" json:"blocker_note"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (GroundTask) TableName() string { return "ground_task" }
