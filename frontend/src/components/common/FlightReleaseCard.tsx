import { StatusBadge } from "./StatusBadge";
import { ReleaseButton } from "./ReleaseButton";
import { ReleaseBlockersPanel } from "./ReleaseBlockersPanel";
import { TurnaroundTimeline } from "./TurnaroundTimeline";
import { useTurnaroundProgress } from "../../hooks/useTurnaroundProgress";
import { useFlightTurnaroundStore } from "../../stores/FlightTurnaroundStore";
import { useGroundTaskStore } from "../../stores/GroundTaskStore";
import { useDelayEventStore } from "../../stores/DelayEventStore";
import { useResourceBookingStore } from "../../stores/ResourceBookingStore";
import { useGroundResourceStore } from "../../stores/GroundResourceStore";
import { formatDate } from "../../utils/formatters";
import type { FlightTurnaroundDetail } from "../../types/FlightTurnaround";

type Props = { flight: FlightTurnaroundDetail; expanded?: boolean };

export function FlightReleaseCard({ flight, expanded = false }: Props) {
  const tasks = useGroundTaskStore((state) => state.rows);
  const delays = useDelayEventStore((state) => state.rows);
  const bookings = useResourceBookingStore((state) => state.rows);
  const resources = useGroundResourceStore((state) => state.rows);
  const serverBlockers = useFlightTurnaroundStore((state) => state.blockersByFlight[flight.id]);
  const progress = useTurnaroundProgress({ flight, tasks, delays, bookings, resources });
  const blockers = serverBlockers ?? progress.blockers;

  return <article className="flight-card panel">
    <header className="flight-card-head">
      <div>
        <p className="eyebrow">{flight.stand_no} · {flight.aircraft_reg}</p>
        <h2>{flight.flight_no}</h2>
        <small>{formatDate(flight.arrival_time)} → {formatDate(flight.departure_time)}</small>
      </div>
      <StatusBadge value={flight.turnaround_status} />
    </header>
    <div className="flight-metrics">
      <span>任务签收 <b>{progress.taskCompletionRate}%</b></span>
      <span>未关闭延误 <b>{blockers.delays.length}</b></span>
      <span>ACTIVE 预约 <b>{progress.activeBookingCount}</b></span>
      <span>冲突预约 <b>{blockers.bookings.length}</b></span>
    </div>
    {expanded ? <TurnaroundTimeline flight={flight} tasks={tasks} bookings={bookings} /> : null}
    <ReleaseBlockersPanel blockers={blockers} compact={!expanded} />
    <footer className="flight-card-foot">
      <span>{flight.ready_at ? `READY @ ${formatDate(flight.ready_at)}` : "放行为原子操作：失败不改动航班、任务和预约"}</span>
      <ReleaseButton flightId={flight.id} status={flight.turnaround_status} ready={progress.ready} />
    </footer>
  </article>;
}
