import { useEffect } from "react";

import { useDelayEventStore } from "../stores/DelayEventStore";
import { useFlightTurnaroundStore } from "../stores/FlightTurnaroundStore";
import { useReleaseStore } from "../stores/ReleaseStore";
import { DelayTag } from "../components/common/DelayTag";
import { formatDate } from "../utils/formatters";

export function DelaysPage() {
  const delays = useDelayEventStore((state) => state.rows);
  const loadDelays = useDelayEventStore((state) => state.load);
  const turnarounds = useFlightTurnaroundStore((state) => state.rows);
  const loadTurnarounds = useFlightTurnaroundStore((state) => state.load);
  const resolveDelay = useReleaseStore((state) => state.resolveDelay);
  const busyIds = useReleaseStore((state) => state.resolvingDelayIds);

  useEffect(() => {
    void Promise.all([loadDelays(), loadTurnarounds()]);
  }, [loadDelays, loadTurnarounds]);

  const totalMinutes = delays.reduce((sum, item) => sum + item.minutes, 0);
  const openMinutes = delays.filter((item) => !item.resolved_at).reduce((sum, item) => sum + item.minutes, 0);
  const flightOf = (id: number) => turnarounds.find((item) => item.id === id);

  return (
    <section className="page-list">
      <header className="page-header">
        <div>
          <p className="eyebrow">ground-turn</p>
          <h1>延误归因</h1>
          <p className="subtitle">未关闭延误必须归零，否则航班整次拒绝放行</p>
        </div>
      </header>

      <section className="metrics">
        <div className="stat"><span>延误事件</span><strong>{delays.length}</strong></div>
        <div className="stat"><span>累计影响</span><strong>{totalMinutes} 分钟</strong></div>
        <div className="stat"><span>未关闭影响</span><strong>{openMinutes} 分钟</strong></div>
      </section>

      <div className="cards">
        {delays.map((delay) => {
          const open = !delay.resolved_at;
          return (
            <article key={delay.id} className={`panel delay-card ${open ? "open" : "closed"}`}>
              <div className="flight-row-head">
                <h3>#{delay.id} <small>{flightOf(delay.turnaround_id)?.flight_no ?? `航班 #${delay.turnaround_id}`}</small></h3>
                <DelayTag value={open ? "OPEN" : "RESOLVED"} />
              </div>
              <p>{delay.delay_type} · +{delay.minutes} 分钟 · 责任团队 {delay.responsibility_team}</p>
              <p className="note">{delay.root_cause}</p>
              <p className="hint">{open ? "未关闭" : `已关闭于 ${formatDate(delay.resolved_at as string)}`}</p>
              {open && (
                <button
                  type="button"
                  className="mini-btn"
                  disabled={Boolean(busyIds[delay.id])}
                  onClick={() => void resolveDelay(delay.id)}
                >
                  {busyIds[delay.id] ? "处理中…" : "关闭延误"}
                </button>
              )}
            </article>
          );
        })}
      </div>
    </section>
  );
}
