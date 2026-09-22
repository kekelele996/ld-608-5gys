export const ERROR_MESSAGES = {
  AUTH_REQUIRED: "请先登录后再继续操作",
  RBAC_DENIED: "当前角色没有执行该动作的权限",
  VALIDATION_FAILED: "表单字段缺失或格式错误",
  RATE_LIMITED: "请求过于频繁，请稍后再试",
  TURNAROUND_NOT_FOUND: "航班过站不存在",
  TURNAROUND_ALREADY_RELEASED: "航班已经放行，请勿重复提交",
  RELEASE_BLOCKED: "放行条件未全部满足",
  CONCURRENT_RELEASE: "放行请求正在处理或已被其他请求完成"
};
