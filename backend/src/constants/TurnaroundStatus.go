package constants

// TurnaroundStatus 枚举：航班过站状态。
// ARRIVING 待入位 / ON_STAND 已靠桥 / IN_SERVICE 保障中 / READY 已放行 / DEPARTED 已离港 / DELAYED 延误挂起。
// 出现位置：models、types、constructors、services、controllers、logTemplates、errorMessages、筛选器、展示组件。
var TurnaroundStatus = []string{"ARRIVING", "ON_STAND", "IN_SERVICE", "READY", "DEPARTED", "DELAYED"}

const (
	TurnaroundArriving  = "ARRIVING"
	TurnaroundOnStand   = "ON_STAND"
	TurnaroundInService = "IN_SERVICE"
	TurnaroundReady     = "READY"
	TurnaroundDeparted  = "DEPARTED"
	TurnaroundDelayed   = "DELAYED"
)

// ReadyTerminalStatuses 已完成放行的状态，重复放行直接拒绝。
var ReadyTerminalStatuses = map[string]bool{
	TurnaroundReady:    true,
	TurnaroundDeparted: true,
}
