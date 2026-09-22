// LOG_TEMPLATES 与后端 constants/logTemplates.go 对应；写操作均落审计日志。
export const LOG_TEMPLATES = {
  FlightTurnaround: [
    "航班过站创建：%s 机位 %s",
    "航班过站更新：#%d 航班 %s",
    "航班过站状态变更：#%d %s -> %s",
    "航班放行联动：#%d 航班 %s 进入 READY，释放 ACTIVE 预约 %d 份，签收任务 %d 项"
  ],
  GroundTask: [
    "地勤任务创建：#%d %s 派发班组 #%d",
    "地勤任务更新：#%d %s",
    "地勤任务状态变更：#%d %s -> %s",
    "地勤任务签收：#%d %s 航班 #%d 任务与过站时间线同步刷新"
  ],
  GroundResource: ["保障资源创建：%s %s", "保障资源更新：#%d %s", "保障资源状态变更：#%d %s -> %s", "保障资源导出：#%d %s"],
  ResourceBooking: [
    "资源预约创建：#%d 资源 #%d 航班 #%d",
    "资源预约更新：#%d 状态 %s",
    "资源预约状态变更：#%d %s -> %s",
    "资源预约批量释放：航班 #%d 一次性释放 ACTIVE 预约 %d 份，资源回到 AVAILABLE"
  ],
  DelayEvent: [
    "延误事件创建：#%d 航班 #%d +%d 分钟",
    "延误事件更新：#%d %s",
    "延误事件归因：#%d 责任团队 %s",
    "延误事件关闭：#%d 航班 #%d 未关闭延误归零"
  ]
} as const;
