import type { ReleaseBlockers } from "../../types/Release";

type Props = {
  blockers: ReleaseBlockers | undefined;
  compact?: boolean;
};

export function ReleaseBlockersPanel({ blockers, compact = false }: Props) {
  if (!blockers || (blockers.tasks.length === 0 && blockers.delays.length === 0 && blockers.bookings.length === 0)) {
    return <div className="empty">暂无阻塞明细；放行时 ACTIVE 预约会随 READY 一次性释放。</div>;
  }

  return <div className="blocker-panel">
    <div className="blocker-summary">
      <b>{blockers.tasks.length + blockers.delays.length + blockers.bookings.length}</b>
      <span>项阻塞</span>
    </div>
    <div className={compact ? "blocker-groups compact" : "blocker-groups"}>
      <section>
        <h3>任务阻塞 · {blockers.tasks.length}</h3>
        {blockers.tasks.map((task) => <article key={`task-${task.id}`}>
          <strong>#{task.id} {task.task_type}</strong>
          <span>{task.reason} · 当前 {task.status}</span>
          {task.blocker_note ? <small>{task.blocker_note}</small> : null}
        </article>)}
      </section>
      <section>
        <h3>延误阻塞 · {blockers.delays.length}</h3>
        {blockers.delays.map((delay) => <article key={`delay-${delay.id}`}>
          <strong>#{delay.id} {delay.delay_type} · {delay.minutes} 分钟</strong>
          <span>{delay.reason}</span>
          <small>{delay.root_cause} / {delay.responsibility_team}</small>
        </article>)}
      </section>
      <section>
        <h3>预约阻塞 · {blockers.bookings.length}</h3>
        {blockers.bookings.map((booking) => <article key={`booking-${booking.id}`}>
          <strong>#{booking.id} {booking.resource_code}</strong>
          <span>{booking.reason} · 当前 {booking.booking_status}</span>
          {booking.conflict_reason ? <small>{booking.conflict_reason}</small> : null}
        </article>)}
      </section>
    </div>
  </div>;
}
