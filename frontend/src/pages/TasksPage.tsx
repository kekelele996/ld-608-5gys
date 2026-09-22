import { useEffect } from "react";

import { useGroundTaskStore } from "../stores/GroundTaskStore";
import { useFlightTurnaroundStore } from "../stores/FlightTurnaroundStore";
import { useReleaseStore } from "../stores/ReleaseStore";
import { StatusBadge } from "../components/common/StatusBadge";
import { GroundTaskTypeText } from "../constants/GroundTaskType";
import { GroundTaskStatusText } from "../constants/GroundTaskStatus";
import { formatDate } from "../utils/formatters";

export function TasksPage() {
  const tasks = useGroundTaskStore((state) => state.rows);
  const loadTasks = useGroundTaskStore((state) => state.load);
  const turnarounds = useFlightTurnaroundStore((state) => state.rows);
  const loadTurnarounds = useFlightTurnaroundStore((state) => state.load);
  const signTask = useReleaseStore((state) => state.signTask);
  const busyIds = useReleaseStore((state) => state.signingTaskIds);

  useEffect(() => {
    void Promise.all([loadTasks(), loadTurnarounds()]);
  }, [loadTasks, loadTurnarounds]);

  const flightOf = (id: number) => turnarounds.find((item) => item.id === id);

  return (
    <section className="page-list">
      <header className="page-header">
        <div>
          <p className="eyebrow">ground-turn</p>
          <h1>地勤任务</h1>
          <p className="subtitle">所有任务签收（SIGNED）后航班才允许放行</p>
        </div>
      </header>
      <div className="table-panel">
        <table>
          <thead>
            <tr>
              <th>#</th><th>航班</th><th>任务</th><th>班组</th><th>截止</th><th>状态</th><th>阻塞备注</th><th></th>
            </tr>
          </thead>
          <tbody>
            {tasks.map((task) => {
              const flight = flightOf(task.turnaround_id);
              const signed = task.status === "SIGNED";
              return (
                <tr key={task.id} className={signed ? "row-signed" : ""}>
                  <td>{task.id}</td>
                  <td>{flight?.flight_no ?? `#${task.turnaround_id}`}</td>
                  <td>{GroundTaskTypeText[task.task_type as keyof typeof GroundTaskTypeText] ?? task.task_type}</td>
                  <td>#{task.team_id}</td>
                  <td>{formatDate(task.deadline)}</td>
                  <td><StatusBadge value={task.status} /></td>
                  <td className="note">{task.blocker_note || "—"}</td>
                  <td>
                    <button
                      type="button"
                      className="mini-btn"
                      disabled={signed || Boolean(busyIds[task.id])}
                      onClick={() => void signTask(task.id)}
                    >
                      {signed ? GroundTaskStatusText.SIGNED : busyIds[task.id] ? "处理中…" : "签收"}
                    </button>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </section>
  );
}
