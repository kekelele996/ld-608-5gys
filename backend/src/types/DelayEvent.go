package types

import "time"

type DelayEventDTO struct {
	ID                 int64      `json:"id"`
	TurnaroundID       int64      `json:"turnaround_id"`
	DelayType          string     `json:"delay_type"`
	Minutes            int        `json:"minutes"`
	RootCause          string     `json:"root_cause"`
	ResponsibilityTeam string     `json:"responsibility_team"`
	ResolvedAt         *time.Time `json:"resolved_at"`
}

// ResolveDelayRequest 关闭延误事件；关闭后未关闭延误计数归零。
type ResolveDelayRequest struct {
	Actor string `json:"actor"`
}
