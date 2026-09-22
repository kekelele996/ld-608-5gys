import { STATUS_TEXT } from "../../constants/statusText";

// StatusBadge 状态徽章：优先使用中文文案映射，未知状态回退展示原始值。
export function StatusBadge({ value }: { value: string }) {
  const maps = [
    STATUS_TEXT.TurnaroundStatus,
    STATUS_TEXT.GroundTaskStatus,
    STATUS_TEXT.BookingStatus,
    STATUS_TEXT.ResourceStatus,
    STATUS_TEXT.GroundTaskType
  ] as const;
  const text = maps.map((entry) => (entry as Record<string, string>)[value]).find(Boolean) ?? String(value).replace(/_/g, " ");
  return (
    <span className={"badge " + String(value).toLowerCase().replace(/_/g, "-")} title={value}>
      {text}
    </span>
  );
}
