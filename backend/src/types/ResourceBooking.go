package types

import "time"

type ResourceBookingDTO struct {
	ID             int64      `json:"id"`
	ResourceID     int64      `json:"resource_id"`
	ResourceCode   string     `json:"resource_code"`
	TurnaroundID   int64      `json:"turnaround_id"`
	TaskID         int64      `json:"task_id"`
	StartTime      time.Time  `json:"start_time"`
	EndTime        time.Time  `json:"end_time"`
	BookingStatus  string     `json:"booking_status"`
	ConflictReason string     `json:"conflict_reason"`
	ReleasedAt     *time.Time `json:"released_at"`
}

// ResolveBookingRequest 人工解决预约冲突（CONFLICT 硬阻塞），解决后回到 ACTIVE，由放行事务统一释放。
type ResolveBookingRequest struct {
	Actor string `json:"actor"`
}
