package constants

// ResourceStatus 枚举：保障资源台账状态。
// AVAILABLE 可用 / BOOKED 被预约占用 / MAINTENANCE 维护中 / OFFLINE 下线。
var ResourceStatus = []string{"AVAILABLE", "BOOKED", "MAINTENANCE", "OFFLINE"}

const (
	ResourceAvailable   = "AVAILABLE"
	ResourceBooked      = "BOOKED"
	ResourceMaintenance = "MAINTENANCE"
	ResourceOffline     = "OFFLINE"
)
