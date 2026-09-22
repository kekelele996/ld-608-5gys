export const GroundTaskStatus = ["DISPATCHED", "SIGNED", "BLOCKED"] as const;
export type GroundTaskStatus = (typeof GroundTaskStatus)[number];
export const GroundTaskStatusText: Record<GroundTaskStatus, string> = {
  DISPATCHED: "已派发",
  SIGNED: "已签收",
  BLOCKED: "阻塞中"
};
