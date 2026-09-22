export const ERROR_MESSAGES = {
  AUTH_REQUIRED: "请先登录后再继续操作",
  RBAC_DENIED: "当前角色没有执行该动作的权限",
  VALIDATION_FAILED: "表单字段缺失或格式错误",
  RATE_LIMITED: "请求过于频繁，请稍后再试",
  TURNAROUND_NOT_FOUND: "航班过站不存在",
  RELEASE_BLOCKED: "放行被阻塞，请先清空任务、延误与预约阻塞清单",
  RELEASE_ALREADY_DONE: "航班已完成放行，请勿重复提交",
  RELEASE_RACE_LOST: "并发提交仅生效一次，本次请求已忽略",
  TASK_NOT_FOUND: "地勤任务不存在",
  TASK_ALREADY_SIGNED: "地勤任务已签收，请勿重复签收",
  BOOKING_NOT_FOUND: "资源预约不存在",
  DELAY_NOT_FOUND: "延误事件不存在",
  INTERNAL_ERROR: "服务内部错误"
} as const;
