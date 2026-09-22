import { useEffect } from "react";
import { FlightReleaseCard } from "../components/common/FlightReleaseCard";
import { useFlightTurnaroundStore } from "../stores/FlightTurnaroundStore";
import { useGroundTaskStore } from "../stores/GroundTaskStore";
import { useDelayEventStore } from "../stores/DelayEventStore";
import { useResourceBookingStore } from "../stores/ResourceBookingStore";
import { useGroundResourceStore } from "../stores/GroundResourceStore";

export function TurnaroundsPage() {
  const flights = useFlightTurnaroundStore((state) => state.rows);
  const notice = useFlightTurnaroundStore((state) => state.notice);
  const clearNotice = useFlightTurnaroundStore((state) => state.clearNotice);
  const loadFlights = useFlightTurnaroundStore((state) => state.load);
  const loadTasks = useGroundTaskStore((state) => state.load);
  const loadDelays = useDelayEventStore((state) => state.load);
  const loadBookings = useResourceBookingStore((state) => state.load);
  const loadResources = useGroundResourceStore((state) => state.load);

  useEffect(() => {
    void Promise.all([loadFlights(), loadTasks(), loadDelays(), loadBookings(), loadResources()]);
  }, [loadFlights, loadTasks, loadDelays, loadBookings, loadResources]);

  return <section className="page-grid">
    <div className="page-head">
      <div>
        <p className="eyebrow">turnaround timeline</p>
        <h1>航班过站</h1>
      </div>
      {notice ? <button className={`notice ${notice.type}`} onClick={clearNotice}>{notice.message}</button> : null}
    </div>
    <section className="flight-list detailed">
      {flights.map((flight) => <FlightReleaseCard key={flight.id} flight={flight} expanded />)}
    </section>
  </section>;
}
