import { useEffect } from "react";
import { StatCard } from "../components/common/StatCard";
import { FlightReleaseCard } from "../components/common/FlightReleaseCard";
import { useFlightTurnaroundStore } from "../stores/FlightTurnaroundStore";
import { useGroundTaskStore } from "../stores/GroundTaskStore";
import { useDelayEventStore } from "../stores/DelayEventStore";
import { useResourceBookingStore } from "../stores/ResourceBookingStore";
import { useGroundResourceStore } from "../stores/GroundResourceStore";

export function DashboardPage() {
  const flights = useFlightTurnaroundStore((state) => state.rows);
  const notice = useFlightTurnaroundStore((state) => state.notice);
  const clearNotice = useFlightTurnaroundStore((state) => state.clearNotice);
  const loadFlights = useFlightTurnaroundStore((state) => state.load);
  const loadTasks = useGroundTaskStore((state) => state.load);
  const loadDelays = useDelayEventStore((state) => state.load);
  const loadBookings = useResourceBookingStore((state) => state.load);
  const loadResources = useGroundResourceStore((state) => state.load);

  const bookings = useResourceBookingStore((state) => state.rows);

  useEffect(() => {
    void Promise.all([loadFlights(), loadTasks(), loadDelays(), loadBookings(), loadResources()]);
  }, [loadFlights, loadTasks, loadDelays, loadBookings, loadResources]);

  const blockedFlights = flights.filter((flight) => flight.turnaround_status !== "READY" && flight.turnaround_status !== "DEPARTED");
  const readyFlights = flights.filter((flight) => flight.turnaround_status === "READY");
  const activeBookings = bookings.filter((booking) => booking.booking_status === "ACTIVE").length;

  return <section className="page-grid">
    <div className="page-head">
      <div>
        <p className="eyebrow">release control</p>
        <h1>过站运行看板</h1>
      </div>
      {notice ? <button className={`notice ${notice.type}`} onClick={clearNotice}>{notice.message}</button> : null}
    </div>
    <section className="metrics">
      <StatCard label="待放行航班" value={blockedFlights.length} />
      <StatCard label="READY 航班" value={readyFlights.length} />
      <StatCard label="ACTIVE 预约" value={activeBookings} />
    </section>
    <section className="flight-list">
      {flights.map((flight) => <FlightReleaseCard key={flight.id} flight={flight} />)}
    </section>
  </section>;
}
