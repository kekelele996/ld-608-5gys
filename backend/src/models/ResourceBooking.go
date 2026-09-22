package models

type ResourceBooking struct {
	ID             int     `json:"id"`
	ResourceID     int     `json:"resource_id"`
	TurnaroundID   int     `json:"turnaround_id"`
	TaskID         int     `json:"task_id"`
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	BookingStatus  string  `json:"booking_status"`
	ConflictReason string  `json:"conflict_reason"`
	ReleasedAt     *string `json:"released_at,omitempty"`
}
