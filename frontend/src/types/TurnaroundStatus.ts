export const TurnaroundStatus = ["ARRIVING","ON_STAND","IN_SERVICE","READY","DEPARTED","DELAYED"] as const;
export type TurnaroundStatus = (typeof TurnaroundStatus)[number];
export const TurnaroundStatusText: Record<TurnaroundStatus, string> = {
  ARRIVING: "待入位",
  ON_STAND: "已靠桥",
  IN_SERVICE: "保障中",
  READY: "已放行",
  DEPARTED: "已离港",
  DELAYED: "延误挂起"
};
