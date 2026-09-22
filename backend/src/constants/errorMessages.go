package constants

const (
	AuthRequiredMessage              = "missing token"
	TurnaroundNotFoundMessage        = "航班过站不存在"
	TurnaroundAlreadyReleasedMessage = "航班已经放行，请勿重复提交"
	ReleaseBlockedMessage            = "放行条件未全部满足"
	ConcurrentReleaseMessage         = "放行请求正在处理或已被其他请求完成"
)
