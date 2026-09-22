import { useEffect, useMemo, useState } from "react";

import { useFlightTurnaroundStore } from "../stores/FlightTurnaroundStore";
import { useGroundResourceStore } from "../stores/GroundResourceStore";
import { useReleaseStore } from "../stores/ReleaseStore";
import { useTurnaroundProgress } from "../hooks/useTurnaroundProgress";
import { useResourceConflict } from "../hooks/useResourceConflict";
import { useResourceBookingStore } from "../stores/ResourceBookingStore";
import { useDelayEventStore } from "../stores/DelayEventStore";
import { useGroundTaskStore } from "../stores/GroundTaskStore";
import { StatusBadge } from "../components/common/StatusBadge";
import { StatCard } from "../components/common/StatCard";
import { BlockerDetails } from "../components/common/BlockerDetails";
import { TurnaroundTimeline } from "../components/common/TurnaroundTimeline";
import type { FlightTurnaround } from "../types/FlightTurnaround";
import type { ApiError } from "../types/Api";

export function DashboardPage() {
  const turnarounds = useFlightTurnaroundStore((state) => state.rows);
  const loadTurnarounds = useFlightTurnaroundStore((state) => state.load);
  const loading = useFlightTurnaroundStore((state) => state.loading);
  const tasks = useGroundTaskStore((state) => state.rows);
  const delays = useDelayEventStore((state) => state.rows);
  const bookings = useResourceBookingStore((state) => state.rows);
  const resources = useGroundResourceStore((state) => state.rows);
  const loadTasks = useGroundTaskStore((state) => state.load);
  const loadDelays = useDelayEventStore((state) => state.load);
  const loadBookings = useResourceBookingStore((state) => state.load);
  const loadResources = useGroundResourceStore((state) => state.load);

  const release = useReleaseStore((state) => state.release);
  const signTask = useReleaseStore((state) => state.signTask);
  const resolveDelay = useReleaseStore((state) => state.resolveDelay);
  const resolveBooking = useReleaseStore((state) => state.resolveBooking);
  const releasingIds = useReleaseStore((state) => state.releasingIds);
  const signingTaskIds = useReleaseStore((state) => state.signingTaskIds);
  const resolvingDelayIds = useReleaseStore((state) => state.resolvingDelayIds);
  const resolvingBookingIds = useReleaseStore((state) => state.resolvingBookingIds);
  const lastBlocked = useReleaseStore((state) => state.lastBlocked);
  const notice = useReleaseStore((state) => state.notice);
  const clearBlocked = useReleaseStore((state) => state.clearBlocked);
  const clearNotice = useReleaseStore((state) => state.clearNotice);
  const refreshAll = useReleaseStore((state) => state.refreshAll);

  const [rejectedError, setRejectedError] = useState<Record<number, ApiError | null>>({});

  useEffect(() => {
    void Promise.all([
      loadTurnarounds(),
      loadTasks(),
      loadDelays(),
      loadBookings(),
      loadResources()
    ]);
  }, [loadTurnarounds, loadTasks, loadDelays, loadBookings, loadResources]);

  const stats = useMemo(() => {
    const ready = turnarounds.filter((item) => item.turnaround_status === "READY").length;
    const blocked = turnarounds.filter((item) => !item.gate.releasable && item.turnaround_status !== "READY").length;
    const openDelays = delays.filter((item) => !item.resolved_at).length;
    return { total: turnarounds.length, ready, blocked, openDelays };
  }, [turnarounds, delays]);

  const conflicts = useResourceConflict(bookings);

  async function handleRelease(turnaround: FlightTurnaround) {
    setRejectedError((prev) => ({ ...prev, [turnaround.id]: null }));
    try {
      await release(turnaround.id);
    } catch (error) {
      // RELEASE_BLOCKED / RELEASE_ALREADY_DONE / RELEASE_RACE_LOST 都展示在卡片上。
      setRejectedError((prev) => ({ ...prev, [turnaround.id]: error as ApiError }));
    }
  }

  return (
    <section className="dashboard">
      <header className="page-header">
        <div>
          <p className="eyebrow">ground-turn</p>
          <h1>过站运行看板</h1>
          <p className="subtitle">放行前联动校验：任务全 SIGNED · 未关闭延误为零 · 阻塞预约清零</p>
        </div>
        <button type="button" className="primary-btn" onClick={() => void refreshAll()}>刷新看板</button>
      </header>

      <section className="metrics">
        <StatCard label="在港航班" value={stats.total} />
        <StatCard label="已放行 READY" value={stats.ready} />
        <StatCard label="门禁阻塞航班" value={stats.blocked} />
        <StatCard label="未关闭延误" value={stats.openDelays} />
        <StatCard label="ACTIVE 预约" value={conflicts.activeCount} />
        <StatCard label="预约冲突" value={conflicts.conflictCount} />
        <StatCard label="资源占用/可用" value={`${resources.filter((r) => r.availability_status === "BOOKED").length}/${resources.length}`} />
      </section>

      {notice && (
        <div className="banner ok" role="status">
          <span>{notice}</span>
          <button type="button" onClick={clearNotice}>×</button>
        </div>
      )}
      {lastBlocked && (
        <div className="banner warn" role="alert">
          <span>
            航班 #{lastBlocked.turnaroundId} 放行被拒绝：{lastBlocked.error.message}
            {lastBlocked.error.code === "RELEASE_ALREADY_DONE" ? "（重复放行不生效）" : ""}
            {lastBlocked.error.code === "RELEASE_RACE_LOST" ? "（并发提交仅生效一次）" : ""}
          </span>
          <button type="button" onClick={clearBlocked}>×</button>
        </div>
      )}

      {loading && turnarounds.length === 0 && <p className="hint">加载中…</p>}

      <section className="turnaround-grid">
        {turnarounds.map((turnaround) => (
          <TurnaroundCard
            key={turnaround.id}
            turnaround={turnaround}
            totalTasks={tasks.filter((task) => task.turnaround_id === turnaround.id).length}
            busy={{
              releasing: Boolean(releasingIds[turnaround.id]),
              task: signingTaskIds,
              delay: resolvingDelayIds,
              booking: resolvingBookingIds
            }}
            rejected={rejectedError[turnaround.id] ?? null}
            onRelease={() => void handleRelease(turnaround)}
            onSignTask={(id) => void signTask(id)}
            onResolveDelay={(id) => void resolveDelay(id)}
            onResolveBooking={(id) => void resolveBooking(id)}
          />
        ))}
      </section>
    </section>
  );
}

function TurnaroundCard({
  turnaround,
  totalTasks,
  busy,
  rejected,
  onRelease,
  onSignTask,
  onResolveDelay,
  onResolveBooking
}: {
  turnaround: FlightTurnaround;
  totalTasks: number;
  busy: {
    releasing: boolean;
    task: Record<number, boolean>;
    delay: Record<number, boolean>;
    booking: Record<number, boolean>;
  };
  rejected: ApiError | null;
  onRelease: () => void;
  onSignTask: (id: number) => void;
  onResolveDelay: (id: number) => void;
  onResolveBooking: (id: number) => void;
}) {
  const progress = useTurnaroundProgress(turnaround);
  // 放行拒绝时优先使用错误体 details 里的阻塞清单（与后端事务内门禁快照一致）。
  const gate = rejected?.details
    ? {
        releasable: false,
        unsigned_tasks: rejected.details.unsigned_tasks ?? turnaround.gate.unsigned_tasks,
        open_delays: rejected.details.open_delays ?? turnaround.gate.open_delays,
        blocking_bookings: rejected.details.blocking_bookings ?? turnaround.gate.blocking_bookings
      }
    : turnaround.gate;

  return (
    <article className={`flight-card ${progress.isReady ? "ready" : ""} ${!gate.releasable && !progress.isReady ? "blocked" : ""}`}>
      <header className="flight-head">
        <div>
          <h3>{turnaround.flight_no} <small>{turnaround.aircraft_reg} · 机位 {turnaround.stand_no}</small></h3>
          <div className="progress-line">
            <div className="progress-bar" style={{ width: `${progress.percent}%` }} />
          </div>
          <small>任务签收 {progress.signed}/{progress.total || totalTasks}（{progress.percent}%）</small>
        </div>
        <StatusBadge value={turnaround.turnaround_status} />
      </header>

      <TurnaroundTimeline turnaround={turnaround} />

      {turnaround.delay_reason && <p className="delay-reason">延误说明：{turnaround.delay_reason}</p>}

      <BlockerDetails
        gate={gate}
        busy={{ task: busy.task, delay: busy.delay, booking: busy.booking }}
        onSignTask={onSignTask}
        onResolveDelay={onResolveDelay}
        onResolveBooking={onResolveBooking}
      />

      {rejected && (
        <p className="reject-msg">
          {rejected.code === "RELEASE_ALREADY_DONE" && "⚠ 重复放行未生效："}
          {rejected.code === "RELEASE_RACE_LOST" && "⚠ 并发落败："}
          {rejected.code === "RELEASE_BLOCKED" && "⚠ 整次拒绝（数据保持原样）："}
          {rejected.message}
        </p>
      )}

      <button
        type="button"
        className="release-btn"
        disabled={busy.releasing || progress.isReady}
        onClick={onRelease}
      >
        {progress.isReady ? "已放行（READY）" : busy.releasing ? "放行中…" : gate.releasable ? "放行航班" : "尝试放行（存在阻塞）"}
      </button>
    </article>
  );
}
