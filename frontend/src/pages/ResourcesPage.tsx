import { useEffect } from "react";
import { StatusBadge } from "../components/common/StatusBadge";
import { useGroundResourceStore } from "../stores/GroundResourceStore";
import { useResourceBookingStore } from "../stores/ResourceBookingStore";

export function ResourcesPage() {
  const resources = useGroundResourceStore((state) => state.rows);
  const bookings = useResourceBookingStore((state) => state.rows);
  const loadResources = useGroundResourceStore((state) => state.load);
  const loadBookings = useResourceBookingStore((state) => state.load);
  useEffect(() => { void Promise.all([loadResources(), loadBookings()]); }, [loadResources, loadBookings]);

  return <section className="page-grid">
    <div className="page-head"><div><p className="eyebrow">resource scheduling</p><h1>资源调度</h1></div></div>
    <div className="panel table-panel">
      {resources.map((resource) => {
        const resourceBookings = bookings.filter((booking) => booking.resource_id === resource.id);
        return <article className="data-row" key={resource.id}>
          <strong>{resource.resource_code}</strong>
          <span>{resource.location} · {resource.owner_team} · {resourceBookings.length} 条预约</span>
          <small>{resourceBookings.some((booking) => booking.booking_status === "CONFLICT") ? "存在冲突，放行将拒绝" : "ACTIVE 预约会在放行时释放"}</small>
          <StatusBadge value={resource.availability_status} />
        </article>;
      })}
    </div>
  </section>;
}
