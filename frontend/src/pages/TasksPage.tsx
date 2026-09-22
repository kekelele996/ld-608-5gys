import { useEffect } from "react";
import { StatusBadge } from "../components/common/StatusBadge";
import { useGroundTaskStore } from "../stores/GroundTaskStore";

export function TasksPage() {
  const tasks = useGroundTaskStore((state) => state.rows);
  const load = useGroundTaskStore((state) => state.load);
  useEffect(() => { void load(); }, [load]);

  return <section className="page-grid">
    <div className="page-head"><div><p className="eyebrow">ground tasks</p><h1>地勤任务</h1></div></div>
    <div className="panel table-panel">
      {tasks.map((task) => <article className="data-row" key={task.id}>
        <strong>#{task.id} {task.task_type}</strong>
        <span>航班 {task.turnaround_id} · 班组 {task.team_id}</span>
        <small>{task.blocker_note || "无阻塞备注"}</small>
        <StatusBadge value={task.status} />
      </article>)}
    </div>
  </section>;
}
