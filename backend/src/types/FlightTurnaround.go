package types

import "time"

// FlightTurnaroundDTO 列表/详情响应。
type FlightTurnaroundDTO struct {
	ID               int64                 `json:"id"`
	FlightNo         string                `json:"flight_no"`
	AircraftReg      string                `json:"aircraft_reg"`
	StandNo          string                `json:"stand_no"`
	ArrivalTime      time.Time             `json:"arrival_time"`
	DepartureTime    time.Time             `json:"departure_time"`
	TurnaroundStatus string                `json:"turnaround_status"`
	DelayReason      string                `json:"delay_reason"`
	ReadyAt          *time.Time            `json:"ready_at"`
	Progress         TurnaroundProgressDTO `json:"progress"`
	Gate             ReleaseGateDTO        `json:"gate"`
}

// TurnaroundProgressDTO 任务与过站时间线同步刷新用聚合。
type TurnaroundProgressDTO struct {
	TotalTasks     int `json:"total_tasks"`
	SignedTasks    int `json:"signed_tasks"`
	OpenDelays     int `json:"open_delays"`
	ActiveBookings int `json:"active_bookings"`
}

// TaskBlockerDTO 未 SIGNED 的任务阻塞项。
type TaskBlockerDTO struct {
	ID          int64  `json:"id"`
	TaskType    string `json:"task_type"`
	TeamID      int64  `json:"team_id"`
	Status      string `json:"status"`
	Deadline    string `json:"deadline"`
	BlockerNote string `json:"blocker_note"`
	Reason      string `json:"reason"`
}

// DelayBlockerDTO 未关闭延误阻塞项。
type DelayBlockerDTO struct {
	ID                 int64  `json:"id"`
	DelayType          string `json:"delay_type"`
	Minutes            int    `json:"minutes"`
	ResponsibilityTeam string `json:"responsibility_team"`
	RootCause          string `json:"root_cause"`
	Reason             string `json:"reason"`
}

// BookingBlockerDTO 无法被放行事务自动释放的预约阻塞项（CONFLICT 等）。
// 普通 ACTIVE 预约由放行事务一次性释放，不进入此清单。
type BookingBlockerDTO struct {
	ID             int64  `json:"id"`
	ResourceID     int64  `json:"resource_id"`
	ResourceCode   string `json:"resource_code"`
	TaskID         int64  `json:"task_id"`
	BookingStatus  string `json:"booking_status"`
	ConflictReason string `json:"conflict_reason"`
	Reason         string `json:"reason"`
}

// ReleaseGateDTO 放行门禁快照：看板阻塞明细直接消费。
type ReleaseGateDTO struct {
	Releasable       bool                `json:"releasable"`
	UnsignedTasks    []TaskBlockerDTO    `json:"unsigned_tasks"`
	OpenDelays       []DelayBlockerDTO   `json:"open_delays"`
	BlockingBookings []BookingBlockerDTO `json:"blocking_bookings"`
}

// ReleaseResultDTO 放行成功结果。
type ReleaseResultDTO struct {
	TurnaroundID     int64          `json:"turnaround_id"`
	FlightNo         string         `json:"flight_no"`
	TurnaroundStatus string         `json:"turnaround_status"`
	ReadyAt          string         `json:"ready_at"`
	ReleasedBookings int64          `json:"released_bookings"`
	Gate             ReleaseGateDTO `json:"gate"`
}
