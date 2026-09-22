import { useEffect } from "react";
import { StatusBadge } from "../components/common/StatusBadge";
import { useDelayEventStore } from "../stores/DelayEventStore";

export function DelaysPage() {
  const delays = useDelayEventStore((state) => state.rows);
  const load = useDelayEventStore((state) => state.load);
  useEffect(() => { void load(); }, [load]);

  return <section className="page-grid">
    <div className="page-head"><div><p className="eyebrow">delay attribution</p><h1>延误归因</h1></div></div>
    <div className="panel table-panel">
      {delays.map((delay) => <article className="data-row" key={delay.id}>
        <strong>#{delay.id} {delay.delay_type} · {delay.minutes} 分钟</strong>
        <span>航班 {delay.turnaround_id} · {delay.responsibility_team}</span>
        <small>{delay.root_cause}</small>
        <StatusBadge value={delay.resolved_at ? "RESOLVED" : "OPEN"} />
      </article>)}
    </div>
  </section>;
}
