package constants

// 错误码集中定义：service 与 controller 分别包装异常，禁止单点吞掉。
const (
	AuthRequired       = "AUTH_REQUIRED"
	RBACDenied         = "RBAC_DENIED"
	ValidationFailed   = "VALIDATION_FAILED"
	RateLimited        = "RATE_LIMITED"
	TurnaroundNotFound = "TURNAROUND_NOT_FOUND"
	ReleaseBlocked     = "RELEASE_BLOCKED"
	ReleaseAlreadyDone = "RELEASE_ALREADY_DONE"
	ReleaseRaceLost    = "RELEASE_RACE_LOST"
	TaskNotFound       = "TASK_NOT_FOUND"
	TaskAlreadySigned  = "TASK_ALREADY_SIGNED"
	BookingNotFound    = "BOOKING_NOT_FOUND"
	DelayNotFound      = "DELAY_NOT_FOUND"
	InternalError      = "INTERNAL_ERROR"
)
