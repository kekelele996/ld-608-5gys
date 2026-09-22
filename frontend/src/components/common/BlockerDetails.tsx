import type { ReleaseGate, TaskBlocker, DelayBlocker, BookingBlocker } from "../../types/ReleaseGate";

// BlockerDetails 放行阻塞明细：分别展示任务、延误、预约三类清单，并提供就地解阻动作。
export function BlockerDetails({
  gate,
  busy,
  onSignTask,
  onResolveDelay,
  onResolveBooking
}: {
  gate: ReleaseGate;
  busy: { task: Record<number, boolean>; delay: Record<number, boolean>; booking: Record<number, boolean> };
  onSignTask: (id: number) => void;
  onResolveDelay: (id: number) => void;
  onResolveBooking: (id: number) => void;
}) {
  return (
    <div className="blockers">
      <BlockerSection<TaskBlocker>
        title={`未签收任务 ${gate.unsigned_tasks.length}`}
        empty="任务均已签收（SIGNED）"
        items={gate.unsigned_tasks}
        keyOf={(item) => item.id}
        head={(item) => (
          <>
            <span className="tag">{item.task_type}</span>
            <span className="tag">班组 #{item.team_id}</span>
            <span className="tag danger">{item.status}</span>
          </>
        )}
        body={(item) => (
          <>
            <p>{item.reason}</p>
            <small>截止 {item.deadline}{item.blocker_note ? ` · ${item.blocker_note}` : ""}</small>
          </>
        )}
        actionLabel="签收任务"
        busyOf={(item) => Boolean(busy.task[item.id])}
        onAction={(item) => onSignTask(item.id)}
      />

      <BlockerSection<DelayBlocker>
        title={`未关闭延误 ${gate.open_delays.length}`}
        empty="未关闭延误为零"
        items={gate.open_delays}
        keyOf={(item) => item.id}
        head={(item) => (
          <>
            <span className="tag">{item.delay_type}</span>
            <span className="tag warn">+{item.minutes} 分钟</span>
            <span className="tag">{item.responsibility_team}</span>
          </>
        )}
        body={(item) => (
          <>
            <p>{item.reason}</p>
            <small>{item.root_cause}</small>
          </>
        )}
        actionLabel="关闭延误"
        busyOf={(item) => Boolean(busy.delay[item.id])}
        onAction={(item) => onResolveDelay(item.id)}
      />

      <BlockerSection<BookingBlocker>
        title={`阻塞预约 ${gate.blocking_bookings.length}`}
        empty="无阻塞预约（ACTIVE 预约放行时自动释放）"
        items={gate.blocking_bookings}
        keyOf={(item) => item.id}
        head={(item) => (
          <>
            <span className="tag">{item.resource_code}</span>
            <span className="tag danger">{item.booking_status}</span>
          </>
        )}
        body={(item) => (
          <>
            <p>{item.reason}</p>
            {item.conflict_reason && <small>{item.conflict_reason}</small>}
          </>
        )}
        actionLabel="解决冲突"
        busyOf={(item) => Boolean(busy.booking[item.id])}
        onAction={(item) => onResolveBooking(item.id)}
      />
    </div>
  );
}

function BlockerSection<T>({
  title,
  empty,
  items,
  keyOf,
  head,
  body,
  actionLabel,
  busyOf,
  onAction
}: {
  title: string;
  empty: string;
  items: T[];
  keyOf: (item: T) => number;
  head: (item: T) => React.ReactNode;
  body: (item: T) => React.ReactNode;
  actionLabel: string;
  busyOf: (item: T) => boolean;
  onAction: (item: T) => void;
}) {
  return (
    <div className="blocker-section">
      <h4>{title}</h4>
      {items.length === 0 && <p className="blocker-ok">✓ {empty}</p>}
      {items.map((item) => (
        <div key={keyOf(item)} className="blocker-item">
          <div className="blocker-tags">{head(item)}</div>
          {body(item)}
          <button
            type="button"
            className="mini-btn"
            disabled={busyOf(item)}
            onClick={() => onAction(item)}
          >
            {busyOf(item) ? "处理中…" : actionLabel}
          </button>
        </div>
      ))}
    </div>
  );
}
