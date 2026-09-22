// GroundTaskStatus 地勤任务状态：放行门禁要求全部 SIGNED。
export const GroundTaskStatus = ["DISPATCHED", "IN_PROGRESS", "SIGNED", "BLOCKED"] as const;
export type GroundTaskStatus = (typeof GroundTaskStatus)[number];

export const GroundTaskStatusText: Record<GroundTaskStatus, string> = {
  DISPATCHED: "已派发",
  IN_PROGRESS: "作业中",
  SIGNED: "已签收",
  BLOCKED: "阻塞"
};
