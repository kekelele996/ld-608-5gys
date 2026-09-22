package constants

// GroundTaskStatus 枚举：地勤任务流转状态。
// DISPATCHED 已派发 / IN_PROGRESS 作业中 / SIGNED 已签收完成 / BLOCKED 阻塞挂起。
// 放行门禁要求：航班下所有任务必须为 SIGNED。
const GroundTaskStatus = "GROUND_TASK_STATUS"

var GroundTaskStatuses = []string{"DISPATCHED", "IN_PROGRESS", "SIGNED", "BLOCKED"}

const (
	TaskDispatched = "DISPATCHED"
	TaskInProgress = "IN_PROGRESS"
	TaskSigned     = "SIGNED"
	TaskBlocked    = "BLOCKED"
)
