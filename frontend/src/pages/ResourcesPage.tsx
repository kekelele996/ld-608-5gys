import { useEffect } from "react";

import { useGroundResourceStore } from "../stores/GroundResourceStore";
import { useResourceBookingStore } from "../stores/ResourceBookingStore";
import { useReleaseStore } from "../stores/ReleaseStore";
import { useResourceConflict } from "../hooks/useResourceConflict";
import { StatusBadge } from "../components/common/StatusBadge";
import { BookingStatusText } from "../constants/BookingStatus";
import { formatDate } from "../utils/formatters";

export function ResourcesPage() {
  const resources = useGroundResourceStore((state) => state.rows);
  const loadResources = useGroundResourceStore((state) => state.load);
  const bookings = useResourceBookingStore((state) => state.rows);
  const loadBookings = useResourceBookingStore((state) => state.load);
  const resolveBooking = useReleaseStore((state) => state.resolveBooking);
  const busyIds = useReleaseStore((state) => state.resolvingBookingIds);

  useEffect(() => {
    void Promise.all([loadResources(), loadBookings()]);
  }, [loadResources, loadBookings]);

  const conflicts = useResourceConflict(bookings);

  return (
    <section className="page-list">
      <header className="page-header">
        <div>
          <p className="eyebrow">ground-turn</p>
          <h1>资源调度</h1>
          <p className="subtitle">ACTIVE 预约在航班放行时一次性释放；CONFLICT 必须先人工解决</p>
        </div>
      </header>

      <section className="metrics">
        <StatMini label="ACTIVE" value={conflicts.activeCount} />
        <StatMini label="CONFLICT" value={conflicts.conflictCount} />
        <StatMini label="RELEASED" value={conflicts.releasedCount} />
      </section>

      <div className="split-panels">
        <div className="table-panel">
          <h3>资源台账</h3>
          <table>
            <thead><tr><th>编码</th><th>类型</th><th>位置</th><th>责任队</th><th>状态</th></tr></thead>
            <tbody>
              {resources.map((resource) => (
                <tr key={resource.id}>
                  <td>{resource.resource_code}</td>
                  <td>{resource.resource_type}</td>
                  <td>{resource.location}</td>
                  <td>{resource.owner_team}</td>
                  <td><StatusBadge value={resource.availability_status} /></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <div className="table-panel">
          <h3>预约日历</h3>
          <table>
            <thead><tr><th>#</th><th>资源</th><th>航班</th><th>窗口</th><th>状态</th><th></th></tr></thead>
            <tbody>
              {bookings.map((booking) => (
                <tr key={booking.id} className={booking.booking_status === "CONFLICT" ? "row-conflict" : ""}>
                  <td>{booking.id}</td>
                  <td>{booking.resource_code}</td>
                  <td>#{booking.turnaround_id}</td>
                  <td>{formatDate(booking.start_time)} ~ {formatDate(booking.end_time)}</td>
                  <td>
                    <StatusBadge value={booking.booking_status} />
                    {booking.conflict_reason && <small className="conflict-reason">{booking.conflict_reason}</small>}
                  </td>
                  <td>
                    {booking.booking_status === "CONFLICT" && (
                      <button
                        type="button"
                        className="mini-btn"
                        disabled={Boolean(busyIds[booking.id])}
                        onClick={() => void resolveBooking(booking.id)}
                      >
                        {busyIds[booking.id] ? "处理中…" : "解决冲突"}
                      </button>
                    )}
                    {booking.booking_status === "RELEASED" && (
                      <small>{BookingStatusText.RELEASED} · {booking.released_at ? formatDate(booking.released_at) : ""}</small>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </section>
  );
}

function StatMini({ label, value }: { label: string; value: number }) {
  return <div className="stat"><span>{label}</span><strong>{value}</strong></div>;
}
