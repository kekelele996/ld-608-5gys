package constants

// 错误消息模板：字段/枚举变更时必须同步修改调用处。
const (
	AuthRequiredMessage       = "missing token"
	RBACDeniedMessage         = "当前角色没有执行该动作的权限"
	ValidationFailedMessage   = "表单字段缺失或格式错误"
	RateLimitedMessage        = "请求过于频繁，请稍后再试"
	TurnaroundNotFoundMessage = "航班过站不存在：#%d"
	ReleaseBlockedMessage     = "航班 %s 放行被阻塞，请先清空任务、延误与预约阻塞清单"
	ReleaseAlreadyDoneMessage = "航班 %s 已完成放行（READY），请勿重复提交"
	ReleaseRaceLostMessage    = "航班 %s 的放行并发提交仅生效一次，本次请求已忽略"
	TaskNotFoundMessage       = "地勤任务不存在：#%d"
	TaskAlreadySignedMessage  = "地勤任务 #%d 已签收，请勿重复签收"
	BookingNotFoundMessage    = "资源预约不存在：#%d"
	DelayNotFoundMessage      = "延误事件不存在：#%d"
	InternalErrorMessage      = "服务内部错误"
)
