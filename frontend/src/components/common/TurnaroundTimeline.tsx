import type { FlightTurnaround } from "../../types/FlightTurnaround";
import { formatDate } from "../../utils/formatters";

const MILESTONES = ["ARRIVING", "ON_STAND", "IN_SERVICE", "READY", "DEPARTED"] as const;
const MILESTONE_LABEL: Record<string, string> = {
  ARRIVING: "待入位",
  ON_STAND: "已靠桥",
  IN_SERVICE: "保障中",
  READY: "已放行",
  DEPARTED: "已离港"
};

// TurnaroundTimeline 过站时间线：放行 READY 后与任务同步刷新 readyAt 节点。
export function TurnaroundTimeline({ turnaround }: { turnaround: FlightTurnaround }) {
  const currentIndex = MILESTONES.indexOf(turnaround.turnaround_status as (typeof MILESTONES)[number]);
  const effectiveIndex = currentIndex < 0 ? 2 : currentIndex;

  return (
    <ol className="timeline">
      {MILESTONES.map((milestone, index) => {
        const reached = index <= effectiveIndex;
        const isNow = index === effectiveIndex;
        let timeLabel = "";
        if (milestone === "ARRIVING") timeLabel = formatDate(turnaround.arrival_time);
        if (milestone === "READY") timeLabel = turnaround.ready_at ? formatDate(turnaround.ready_at) : "待放行";
        if (milestone === "DEPARTED") timeLabel = milestone === turnaround.turnaround_status ? formatDate(turnaround.departure_time) : "计划 " + formatDate(turnaround.departure_time);
        return (
          <li key={milestone} className={`timeline-step ${reached ? "reached" : ""} ${isNow ? "now" : ""}`}>
            <span className="dot" />
            <div>
              <strong>{MILESTONE_LABEL[milestone]}</strong>
              {timeLabel && <small>{timeLabel}</small>}
            </div>
          </li>
        );
      })}
    </ol>
  );
}
