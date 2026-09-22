package constants

// BookingStatus 枚举：资源预约状态。
// ACTIVE 生效占用 —— 放行事务通过门禁后一次性释放，不构成硬阻塞
// RELEASED 已释放（资源回到 AVAILABLE）
// CONFLICT 存在冲突未解决 —— 属放行硬阻塞，必须先人工解决，放行事务不会自动释放
const BookingStatus = "BOOKING_STATUS"

var BookingStatuses = []string{"ACTIVE", "RELEASED", "CONFLICT"}

const (
	BookingActive   = "ACTIVE"
	BookingReleased = "RELEASED"
	BookingConflict = "CONFLICT"
)

// GateBlockingBookingStatuses 放行门禁检查中会导致整次拒绝的预约状态。
// 注意：ACTIVE 不在其中 —— ACTIVE 预约是在放行事务内被一次性释放（“通过后预约不再 ACTIVE”）；
// 只有无法安全自动释放的 CONFLICT 才必须在放行前清零。
var GateBlockingBookingStatuses = []string{BookingConflict}
