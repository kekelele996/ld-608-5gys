import type { FlightTurnaroundDetail } from "../../types/FlightTurnaround";
import type { GroundTask } from "../../types/GroundTask";
import type { ResourceBooking } from "../../types/ResourceBooking";
import { StatusBadge } from "./StatusBadge";
import { formatDate } from "../../utils/formatters";

type Props = {
  flight: FlightTurnaroundDetail;
  tasks?: GroundTask[];
  bookings?: ResourceBooking[];
};

export function TurnaroundTimeline({ flight, tasks = [], bookings = [] }: Props) {
  const taskByType = new Map(tasks.map((task) => [task.task_type, task]));
  const bookingByTask = new Map(bookings.map((booking) => [booking.task_id, booking]));

  return <div className="timeline-card">
    <div className="timeline-head">
      <strong>{flight.flight_no} · {flight.aircraft_reg}</strong>
      <StatusBadge value={flight.turnaround_status} />
    </div>
    <ol className="timeline-list">
      <li><span>航班到达</span><time>{formatDate(flight.arrival_time)}</time><StatusBadge value="ARRIVING" /></li>
      {flight.timeline.filter((event) => event.code.startsWith("TASK_")).map((event) => {
        const taskType = event.code.replace("TASK_", "");
        const task = taskByType.get(taskType);
        const booking = task ? bookingByTask.get(task.id) : undefined;
        return <li key={event.code}>
          <span>{event.label} 任务</span>
          <time>{formatDate(event.occurred_at)}</time>
          <StatusBadge value={event.status} />
          {booking ? <StatusBadge value={booking.booking_status} /> : null}
        </li>;
      })}
      <li className={flight.turnaround_status === "READY" ? "done" : "pending"}>
        <span>过站放行</span>
        <time>{flight.ready_at ? formatDate(flight.ready_at) : "等待联动校验"}</time>
        <StatusBadge value={flight.turnaround_status === "READY" ? "READY" : "WAITING"} />
      </li>
    </ol>
  </div>;
}
