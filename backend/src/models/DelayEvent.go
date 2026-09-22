package models

import "time"

// DelayEvent 延误事件。ResolvedAt 为空即“未关闭”，未关闭数必须为零才允许放行。
type DelayEvent struct {
	ID                 int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	TurnaroundID       int64      `gorm:"column:turnaround_id;index" json:"turnaround_id"`
	DelayType          string     `gorm:"column:delay_type;size:32" json:"delay_type"`
	Minutes            int        `gorm:"column:minutes" json:"minutes"`
	RootCause          string     `gorm:"column:root_cause;size:255" json:"root_cause"`
	ResponsibilityTeam string     `gorm:"column:responsibility_team;size:64" json:"responsibility_team"`
	ResolvedAt         *time.Time `gorm:"column:resolved_at" json:"resolved_at"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (DelayEvent) TableName() string { return "delay_event" }
