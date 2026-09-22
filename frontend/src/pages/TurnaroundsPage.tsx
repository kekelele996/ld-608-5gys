import { useEffect } from "react";

import { useFlightTurnaroundStore } from "../stores/FlightTurnaroundStore";
import { useReleaseStore } from "../stores/ReleaseStore";
import { useTurnaroundProgress } from "../hooks/useTurnaroundProgress";
import { StatusBadge } from "../components/common/StatusBadge";
import { TurnaroundTimeline } from "../components/common/TurnaroundTimeline";
import type { FlightTurnaround } from "../types/FlightTurnaround";
import { formatDate } from "../utils/formatters";

export function TurnaroundsPage() {
  const rows = useFlightTurnaroundStore((state) => state.rows);
  const load = useFlightTurnaroundStore((state) => state.load);
  const release = useReleaseStore((state) => state.release);
  const releasingIds = useReleaseStore((state) => state.releasingIds);

  useEffect(() => {
    void load();
  }, [load]);

  return (
    <section className="page-list">
      <header className="page-header">
        <div>
          <p className="eyebrow">ground-turn</p>
          <h1>航班过站</h1>
        </div>
      </header>
      <div className="cards">
        {rows.map((turnaround) => (
          <TurnaroundRow
            key={turnaround.id}
            turnaround={turnaround}
            busy={Boolean(releasingIds[turnaround.id])}
            onRelease={() => void release(turnaround.id).catch(() => undefined)}
          />
        ))}
      </div>
    </section>
  );
}

function TurnaroundRow({ turnaround, busy, onRelease }: {
  turnaround: FlightTurnaround;
  busy: boolean;
  onRelease: () => void;
}) {
  const progress = useTurnaroundProgress(turnaround);
  return (
    <article className="panel flight-row">
      <div className="flight-row-head">
        <h3>{turnaround.flight_no} <small>{turnaround.aircraft_reg} · {turnaround.stand_no}</small></h3>
        <StatusBadge value={turnaround.turnaround_status} />
      </div>
      <p className="hint">
        到港 {formatDate(turnaround.arrival_time)} · 计划离港 {formatDate(turnaround.departure_time)}
        {turnaround.ready_at ? ` · 放行时刻 ${formatDate(turnaround.ready_at)}` : ""}
      </p>
      <TurnaroundTimeline turnaround={turnaround} />
      <div className="gate-summary">
        <span className={progress.unsignedTasks ? "bad" : "good"}>未签收任务 {progress.unsignedTasks}</span>
        <span className={progress.openDelays ? "bad" : "good"}>未关闭延误 {progress.openDelays}</span>
        <span className={progress.blockingBookings ? "bad" : "good"}>阻塞预约 {progress.blockingBookings}</span>
        <span>ACTIVE 预约 {progress.activeBookings}（放行时释放）</span>
      </div>
      <button type="button" className="release-btn" disabled={busy || progress.isReady} onClick={onRelease}>
        {progress.isReady ? "已放行" : busy ? "放行中…" : progress.releasable ? "放行航班" : "存在阻塞，请到看板处理"}
      </button>
    </article>
  );
}
