// BookingStatus 资源预约状态：
// ACTIVE 生效占用（放行通过时由后端一次性释放，不构成硬阻塞）
// RELEASED 已释放（资源回到 AVAILABLE）
// CONFLICT 冲突未解决（放行硬阻塞，必须先人工解决）
export const BookingStatus = ["ACTIVE", "RELEASED", "CONFLICT"] as const;
export type BookingStatus = (typeof BookingStatus)[number];

export const BookingStatusText: Record<BookingStatus, string> = {
  ACTIVE: "占用中",
  RELEASED: "已释放",
  CONFLICT: "冲突"
};
