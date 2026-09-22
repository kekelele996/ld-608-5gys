export const BookingStatus = ["ACTIVE", "RELEASED", "CONFLICT"] as const;
export type BookingStatus = (typeof BookingStatus)[number];
export const BookingStatusText: Record<BookingStatus, string> = {
  ACTIVE: "预约占用",
  RELEASED: "已释放",
  CONFLICT: "冲突未决"
};
